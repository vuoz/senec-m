package db

import (
	"senec-monitor/types"
)

type DbService interface {
	Write(*types.Response2) error
	GetData() ([]types.GetDataDB, error)
	GetSpecificData(int64) (types.GetDataDB, error)
	WriteLocalApiData(types.LocalApiDataWithCorrectTypes) error
	GetLatestFromLocal() (types.LocalApiDataWithCorrectTypesWithTimeStamp, error)
	WriteWeather24Hours(types.WeatherDailyDb) error
}
