package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"senec-monitor/db"
	"senec-monitor/logging"
	pb "senec-monitor/proto"
	"senec-monitor/scheduler"
	"senec-monitor/types"
	"senec-monitor/utils"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn    *websocket.Conn
	retries int
}

type Server struct {
	db              db.DbService
	logger          logging.Logger
	DailyPrediction DailyPrediction

	// this is used to give real time data to clients
	Clients *WsStruct
}
type DailyPrediction struct {
	mu   *sync.Mutex
	Data *[]int32 `json:"data"`
	new  bool
}
type WsStruct struct {
	map_mu     *sync.Mutex
	ws_clients map[uuid.UUID]*Client
}

func NewServer(log logging.Logger, db db.DbService) *Server {
	return &Server{
		db:     db,
		logger: log,
		Clients: &WsStruct{
			map_mu:     &sync.Mutex{},
			ws_clients: make(map[uuid.UUID]*Client),
		},
		DailyPrediction: DailyPrediction{
			mu:   &sync.Mutex{},
			Data: nil,
			new:  false,
		},
	}
}

func (s *Server) Start(c <-chan *types.LocalApiDataWithCorrectTypes, latestWeather *types.LatestWeather, latestTotal *types.LatestTotal, latestLocal *types.LatestLocal, predUrl string, coords types.Cordinate) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", s.wrapHandler(handleIndex))
	serveMux.Handle("/full", s.wrapHandlerWithDB(handleGetData))
	serveMux.Handle("/data", s.wrapHandlerWithDB(handleGetSpecificTs))
	serveMux.Handle("/localLatest", s.wrapHandlerWithLocalLatest(handleGetLocalLatest, latestLocal))
	serveMux.Handle("/subscribe", s.wrapHandlerWithWsMapAndPrediction(handleUpgrade))
	serveMux.Handle("/prediction", s.wrapHandlerWithDailyPrediction(handleGetPrediction))
	s.logger.Info("Started Server")
	if predUrl != "" {
		s.logger.Info("Getting predictions")
		coorsFloat, err := coords.ToFloat()
		if err != nil {
			s.logger.Err("Error converting coordinates to float: ", err)
			return
		}

		retries := 0
		go func() {
			for {
				date := time.Now().Format("2006-01-02")
				pred, err := utils.GetPrediction(predUrl, types.PredictionRequest{Date: date, Coord: coorsFloat})
				if err != nil {
					s.logger.Err("error getting yield prediction: ", err)
					time.Sleep(5 * time.Minute)

					retries++
					if retries > 3 {
						s.logger.Err("Error getting yield prediction. stopping retrieval")
						return
					}
					continue

				}
				retries = 0
				s.DailyPrediction.mu.Lock()

				// we scale it down since the device cant handle this much data
				scaled := make([]int32, len(pred))

				if s.DailyPrediction.Data == nil {
					data := make([]int32, len(pred))
					s.DailyPrediction.Data = &data
				}
				data := *s.DailyPrediction.Data
				for i := 0; i < len(pred); i++ {
					scaled[i] = int32(data[i] * 1000)

				}
				s.DailyPrediction.Data = &scaled
				s.DailyPrediction.mu.Unlock()
				s.Clients.map_mu.Lock()
				pred_proto := pb.Prediction{Prediction: scaled}
				bytes, err := proto.Marshal(&pred_proto)
				if err != nil {
					s.logger.Err("Error marshalling prediction: ", err)
				}
				// send new prediction to all the clients
				for k, v := range s.Clients.ws_clients {
					if err := v.conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
						if v.retries > 3 {
							delete(s.Clients.ws_clients, k)
							continue
						}
						s.Clients.ws_clients[k] = &Client{conn: v.conn, retries: v.retries + 1}
						continue
					}

				}

				s.Clients.map_mu.Unlock()

				timeToSleep := scheduler.ScheduleTo8Am()
				time.Sleep(timeToSleep)
				continue

			}

		}()
	}
	go func() {
		for {
			select {
			case msg := <-c:
				{
					// error doesnt matter since it still returns an emtpy struct
					v, _ := latestTotal.Get()
					latestLocal.Set(*msg)

					var string_data types.LocalApiDataWithCorrectTypesWithTimeStampStringsWithWeather
					if s.DailyPrediction.Data != nil && s.DailyPrediction.new {
						//todo
						string_data = msg.ConvertToStrings(latestWeather.Get(), v, nil)
					} else {
						string_data = msg.ConvertToStrings(latestWeather.Get(), v, nil)
					}

					s.Clients.map_mu.Lock()
					for k, v := range s.Clients.ws_clients {
						if err := v.conn.WriteJSON(string_data); err != nil {
							if v.retries > 3 {
								delete(s.Clients.ws_clients, k)
								continue
							}
							s.Clients.ws_clients[k] = &Client{conn: v.conn, retries: v.retries + 1}
							continue
						}

					}
					s.Clients.map_mu.Unlock()

				}
			default:
				continue

			}

		}

	}()
	if err := http.ListenAndServe("0.0.0.0:6600", serveMux); err != nil {
		s.logger.Info('E', "Error occured starting server: ", err)
	}
}

func (s *Server) wrapHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if err := fn(wr, r); err != nil {
			var handlerError *HandlerError
			ok := errors.As(err, &handlerError)
			if ok {

				// These errors wont be logged, since these are user errors
				wr.WriteHeader(500)
				data, err := json.Marshal(types.HandlerErrorResponse{Error: err.Error()})
				if err != nil {
					return
				}
				wr.Write(data)
				return

			}

			s.logger.Info(err)
			wr.WriteHeader(500)
			wr.Write([]byte("Internal Server Error"))
			return
		}
	}

}
func (s *Server) wrapHandlerWithLocalLatest(fn func(http.ResponseWriter, *http.Request, *types.LatestLocal) error, data *types.LatestLocal) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if err := fn(wr, r, data); err != nil {
			var handlerError *HandlerError
			ok := errors.As(err, &handlerError)
			if ok {

				// These errors wont be logged, since these are user errors
				wr.WriteHeader(500)
				data, err := json.Marshal(types.HandlerErrorResponse{Error: err.Error()})
				if err != nil {
					return
				}
				wr.Write(data)
				return

			}

			s.logger.Info(err)
			wr.WriteHeader(500)
			wr.Write([]byte("Internal Server Error"))
			return
		}
	}

}

func (s *Server) wrapHandlerWithWsMapAndPrediction(fn func(http.ResponseWriter, *http.Request, *WsStruct, *DailyPrediction) error) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if err := fn(wr, r, s.Clients, &s.DailyPrediction); err != nil {

			var handlerError *HandlerError
			ok := errors.As(err, &handlerError)
			if ok {

				// These errors wont be logged, since these are user errors
				data, err := json.Marshal(types.HandlerErrorResponse{Error: err.Error()})
				if err != nil {
					return
				}
				wr.Write(data)
				return

			}

			s.logger.Info(err)
			wr.WriteHeader(500)
			wr.Write([]byte("Internal Server Error"))
			return
		}
	}

}
func (s *Server) wrapHandlerWithDB(fn func(http.ResponseWriter, *http.Request, db.DbService) error) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if err := fn(wr, r, s.db); err != nil {

			var handlerError *HandlerError
			ok := errors.As(err, &handlerError)
			if ok {
				// These errors wont be logged, since these are user errors
				wr.WriteHeader(500)
				data, err := json.Marshal(types.HandlerErrorResponse{Error: err.Error()})
				if err != nil {
					return
				}
				wr.Write(data)
				return

			}
			s.logger.Info(err)
			wr.WriteHeader(500)
			wr.Write([]byte("Internal Server Error"))
			return
		}
	}

}
func (s *Server) wrapHandlerWithDailyPrediction(fn func(http.ResponseWriter, *http.Request, *DailyPrediction) error) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if err := fn(wr, r, &s.DailyPrediction); err != nil {
			var handlerError *HandlerError
			ok := errors.As(err, &handlerError)
			if ok {
				// These errors wont be logged, since these are user errors
				wr.WriteHeader(500)
				data, err := json.Marshal(types.HandlerErrorResponse{Error: err.Error()})
				if err != nil {
					return
				}
				wr.Write(data)
				return

			}
			s.logger.Info(err)
			wr.WriteHeader(500)
			wr.Write([]byte("Internal Server Error"))
			return
		}

	}

}
