package task

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"senec-monitor/db"
	"senec-monitor/logging"
	"senec-monitor/types"
	"senec-monitor/utils"
	"strings"
	"time"
)

type LocalTask struct {
	c   *http.Client
	log logging.Logger
	ip  string
}

func NewLocalTask(log logging.Logger, ip string) *LocalTask {
	// to skip the tls warning
	transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: &transport}
	return &LocalTask{
		c:   c,
		log: log,
		ip:  ip,
	}

}
func CreateAndLoopLocalTask(log logging.Logger, db db.DbService, c chan<- *types.LocalApiDataWithCorrectTypes, ip string) {
	task := NewLocalTask(log, ip)
	task.loop(db, c)
}

func (t *LocalTask) GetData() (types.LocalApiResponse, error) {
	respForError := types.LocalApiResponse{}

	text := `{"ENERGY":{"STAT_STATE":"","GUI_BAT_DATA_POWER":"","GUI_INVERTER_POWER":"","GUI_HOUSE_POW":"","GUI_GRID_POW":"","GUI_BAT_DATA_FUEL_CHARGE":"","GUI_CHARGING_INFO":"","GUI_BOOSTING_INFO":""},"SYS_UPDATE":{"UPDATE_AVAILABLE":""},"STECA":{}}`
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lala.cgi", t.ip), strings.NewReader(text))
	if err != nil {
		return respForError, err
	}
	req.Header = http.Header{
		"X-Requested-With": []string{"XMLHttpRequest"},
		"Content-type":     []string{"application/x-www-form-urlencoded; charset=UTF-8"},
	}
	resp, err := t.c.Do(req)
	if err != nil {
		return respForError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return respForError, fmt.Errorf("status error %d ", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return respForError, err
	}

	var data = &types.LocalApiResponse{}
	if err := json.Unmarshal(body, data); err != nil {
		return respForError, err
	}
	return *data, nil

}
func (t *LocalTask) loop(db db.DbService, c chan<- *types.LocalApiDataWithCorrectTypes) {
	num := 0
	for {
		res, err := t.GetData()
		if err != nil {
			t.log.Err("Error getting data from local api: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		parsedData, err := utils.ParseStringDataToStruct(res)
		if err != nil {
			t.log.Err("Error parsing response to correct types: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		c <- parsedData
		if err := db.WriteLocalApiData(*parsedData); err != nil {
			t.log.Err("Cannot save data to database: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		num++
		if num%10 == 0 {
			t.log.Info("collected local data 10x")
		}

		time.Sleep(10 * time.Second)
	}

}
