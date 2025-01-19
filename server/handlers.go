package server

import (
	"encoding/json"
	"errors"
	"net/http"
	database "senec-monitor/db"
	pb "senec-monitor/proto"
	"senec-monitor/types"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleIndex(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(200)
	w.Write([]byte("Hello"))
	return nil

}

func handleGetData(w http.ResponseWriter, _ *http.Request, db database.DbService) error {

	res, err := db.GetData()
	if err != nil {
		return err
	}
	json, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(json)
	return nil

}
func handleGetSpecificTs(w http.ResponseWriter, r *http.Request, db database.DbService) error {
	if r.URL.Query().Get("ts") == "" {
		return NewHandlerError("Your missing the timestamp!")
	}
	ts := r.URL.Query().Get("ts")
	tsInt, err := strconv.Atoi(ts)
	if err != nil {
		return NewHandlerError("Your timestamp is not correct!")
	}

	res, err := db.GetSpecificData(int64(tsInt))
	if err != nil {
		var userError *database.UserInputError
		ok := errors.As(err, &userError)
		if ok {
			return NewHandlerError(err.Error())
		} else {

			return err
		}
	}
	json, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(json)
	return nil
}
func handleGetLocalLatest(w http.ResponseWriter, _ *http.Request, dataa *types.LatestLocal) error {
	data := dataa.Get()
	json, err := json.Marshal(data)
	if err != nil {
		var userError *database.UserInputError
		ok := errors.As(err, &userError)
		if ok {
			return NewHandlerError(err.Error())

		} else {

			return err
		}
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(json)
	return nil
}
func handleUpgrade(w http.ResponseWriter, r *http.Request, clients *WsStruct, pred *DailyPrediction) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return NewHandlerError("Error opening a websocket connection!")
	}
	id := uuid.New()
	// writing this upon first connection to avoid sending too much data at once
	if pred != nil && pred.Data != nil {
		pred_new := pb.Prediction{Prediction: *pred.Data}
		protoData := pb.Data{Oneof: &pb.Data_Prediction{Prediction: &pred_new}}
		bytes_new, err := proto.Marshal(&protoData)
		if err != nil {
			return err
		}
		if er := conn.WriteMessage(websocket.BinaryMessage, bytes_new); er != nil {
			return err
		}
	}

	clients.map_mu.Lock()
	clients.ws_clients[id] = &Client{retries: 0, conn: conn}
	clients.map_mu.Unlock()
	return nil
}
func handleGetPrediction(w http.ResponseWriter, r *http.Request, predicion *DailyPrediction) error {
	if predicion.Data == nil {
		return json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{Message: "No data yet!"})

	}
	predicion.mu.Lock()
	defer predicion.mu.Unlock()
	if err := json.NewEncoder(w).Encode(predicion); err != nil {
		return err
	}
	return nil
}
