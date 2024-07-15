package server

import (
	"encoding/json"
	"errors"
	"net/http"
	database "senec-monitor/db"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
func handleGetLocalLatest(w http.ResponseWriter, _ *http.Request, db database.DbService) error {
	res, err := db.GetLatestFromLocal()
	if err != nil {
		return err
	}
	json, err := json.Marshal(res)
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
func handleUpgrade(w http.ResponseWriter, r *http.Request, clients *WsStruct) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return NewHandlerError("Error opening a websocket connection!")
	}
	id := uuid.New()
	clients.map_mu.Lock()
	clients.ws_clients[id] = &Client{retries: 0, conn: conn}
	clients.map_mu.Unlock()
	return nil
}
