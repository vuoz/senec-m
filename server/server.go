package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"senec-monitor/db"
	"senec-monitor/logging"
	"senec-monitor/types"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	retries int
}

type Server struct {
	db     db.DbService
	logger logging.Logger

	// this is used to give real time data to clients
	Clients *WsStruct
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
	}
}

func (s *Server) Start(c <-chan *types.LocalApiDataWithCorrectTypes, latestWeather *types.LatestWeather, latestTotal *types.LatestTotal, latestLocal *types.LatestLocal) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", s.wrapHandler(handleIndex))
	serveMux.Handle("/full", s.wrapHandlerWithDB(handleGetData))
	serveMux.Handle("/data", s.wrapHandlerWithDB(handleGetSpecificTs))
	serveMux.Handle("/localLatest", s.wrapHandlerWithLocalLatest(handleGetLocalLatest, latestLocal))
	serveMux.Handle("/subscribe", s.wrapHandlerWithWsMap(handleUpgrade, s.Clients))
	s.logger.Info("Started Server")
	go func() {
		for {
			select {
			case msg := <-c:
				{
					// error doesnt matter since it still returns an emtpy struct
					v, _ := latestTotal.Get()
					latestLocal.Set(*msg)

					string_data := msg.ConvertToStrings(latestWeather.Get(), v)

					for k, v := range s.Clients.ws_clients {
						if err := v.conn.WriteJSON(string_data); err != nil {
							s.Clients.map_mu.Lock()
							if v.retries > 3 {
								delete(s.Clients.ws_clients, k)
								s.Clients.map_mu.Unlock()
								continue
							}
							s.Clients.ws_clients[k] = &Client{conn: v.conn, retries: v.retries + 1}
							s.Clients.map_mu.Unlock()
							continue
						}

					}

				}
			default:
				continue

			}

		}

	}()
	if err := http.ListenAndServe("0.0.0.0:6000", serveMux); err != nil {
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

func (s *Server) wrapHandlerWithWsMap(fn func(http.ResponseWriter, *http.Request, *WsStruct) error, clients *WsStruct) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		if err := fn(wr, r, clients); err != nil {

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
