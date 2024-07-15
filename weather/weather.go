package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"senec-monitor/db"
	"senec-monitor/logging"
	"senec-monitor/scheduler"
	"senec-monitor/types"
	"time"
)

type task struct {
	client *http.Client
	cords  types.Cordinate
}

func (t *task) GetWeatherHourly() (types.WeatherApiResponse, error) {
	dataForErrResp := types.WeatherApiResponse{}
	params := url.Values{}
	params.Set("latitude", t.cords.Lat)
	params.Set("longitude", t.cords.Long)
	params.Add("hourly", "cloud_cover")
	params.Add("hourly", "visibility")
	params.Add("daily", "sunrise")
	params.Add("daily", "sunset")
	params.Add("daily", "sunshine_duration")
	params.Add("daily", "uv_index_max")
	params.Add("daily", "uv_index_clear_sky_max")

	encoded := params.Encode()
	url := "https://api.open-meteo.com/v1/forecast"
	finalUrl := url + "?" + encoded
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return dataForErrResp, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return dataForErrResp, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return dataForErrResp, err
		}
		if string(body) != "" {
			return dataForErrResp, fmt.Errorf("status error %d sending weather request. Error Response: %s", resp.StatusCode, string(body))
		}
		return dataForErrResp,
			fmt.Errorf("status error sending request: %v ", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dataForErrResp, err
	}
	var structuredData types.WeatherApiResponse
	if err := json.Unmarshal(body, &structuredData); err != nil {
		return dataForErrResp, err
	}
	return structuredData, nil
}
func (t *task) GetWeatherDaily() (types.WeatherDailyResp, error) {
	respForError := types.WeatherDailyResp{}
	params := url.Values{}
	params.Set("latitude", t.cords.Lat)

	params.Set("longitude", t.cords.Long)
	params.Add("daily", "weather_code")
	params.Add("daily", "temperature_2m_max")
	params.Add("daily", "temperature_2m_min")
	params.Add("daily", "sunrise")
	params.Add("daily", "sunset")
	params.Add("daily", "daylight_duration")
	params.Add("daily", "sunshine_duration")
	params.Add("daily", "uv_index_max")
	params.Add("daily", "uv_index_clear_sky_max")

	params.Set("timezone", "Europe/Berlin")
	params.Set("forecast_days", "1")

	encoded := params.Encode()
	url := "https://api.open-meteo.com/v1/forecast"
	finalUrl := url + "?" + encoded
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return respForError, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return respForError, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return respForError, err
	}
	defer resp.Body.Close()
	var data types.WeatherDailyResp
	if err := json.Unmarshal(body, &data); err != nil {
		return respForError, err
	}

	return data, nil

}
func newWeatherTask(cords types.Cordinate) *task {
	return &task{
		client: &http.Client{},
		cords:  cords,
	}

}
func GetWeatherEvery10Min(cords types.Cordinate, logger logging.Logger, weatherChan chan<- types.WeatherApiResponse) {
	task := newWeatherTask(cords)
	for {
		res, err := task.GetWeatherHourly()
		if err != nil {
			logger.Log('E', "Error getting weather data ", err)
			time.Sleep(2 * time.Minute)
			continue
		}
		weatherChan <- res
		time.Sleep(10 * time.Minute)

	}

}
func GetWeatherEvery24(cords types.Cordinate, logger logging.Logger, database db.DbService) {
	task := newWeatherTask(cords)
	retry_count := 0

	// for the inital startup
	time.Sleep(scheduler.ScheduleTo8Am())

	for {

		resp, err := task.GetWeatherDaily()
		if err != nil {
			if retry_count >= 5 {
				logger.Log('E', "Hit retry limit on Daily Weather task")
			}
			logger.Log('E', "Error getting Daily Weather: ", err)
			retry_count++
			time.Sleep(10 * time.Second)
			continue

		}
		retry_count = 0

		if err := database.WriteWeather24Hours(resp.ToDbType()); err != nil {
			logger.Log('E', "Error saving weather data ", err)
		}

		// for every reoccuring day
		time.Sleep(scheduler.ScheduleTo8Am())
		// schedule to every day at 6:00

	}

}
func (t *task) GetWeatherHourlyV2() (types.ApiRespHourly, error) {
	respForError := types.ApiRespHourly{}

	params := url.Values{}
	params.Add("latitude", t.cords.Lat)
	params.Add("longitude", t.cords.Long)

	params.Add("hourly", "temperature_2m")
	params.Add("hourly", "rain")
	params.Add("hourly", "showers")
	params.Add("hourly", "cloud_cover")
	params.Add("hourly", "uv_index")
	params.Add("hourly", "uv_index_clear_sky")
	params.Add("daily", "sunset")
	params.Add("daily", "sunrise")

	params.Add("timezone", "Europe/Berlin")
	params.Add("forecast_days", "1")
	params.Add("forecast_hours", "4")
	encoded := params.Encode()
	url := "https://api.open-meteo.com/v1/forecast"
	finalUrl := url + "?" + encoded
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return respForError, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return respForError, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return respForError, err
		}
		if string(body) != "" {
			return respForError, fmt.Errorf("status error %d sending weather request. Error Response: %s", resp.StatusCode, string(body))
		}
		return respForError, fmt.Errorf("status error sending request: %v ", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return respForError, err
	}
	var data types.ApiRespHourly
	if err := json.Unmarshal(body, &data); err != nil {
		return respForError, err
	}

	return data, nil

}

func GetWeatherHourly(cords types.Cordinate, logger logging.Logger, weatherChan chan<- types.ApiRespHourly) {
	t := newWeatherTask(cords)
	retries := 0
	for {
		resp, err := t.GetWeatherHourlyV2()
		if err != nil {
			logger.Log('E', "Error getting weather data: ", err)
			time.Sleep(10 * time.Second)
			retries++
			if retries > 5 {
				break
			}
			continue
		}

		retries = 0
		weatherChan <- resp

		time.Sleep(1 * time.Hour)

	}
	logger.Log('E', "Stopped weather task after 5 retries.")

}
