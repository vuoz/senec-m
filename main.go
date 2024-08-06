package main

import (
	"github.com/joho/godotenv"
	"os"
	"senec-monitor/db"
	"senec-monitor/logging"
	"senec-monitor/server"
	"senec-monitor/task"
	"senec-monitor/types"
	"senec-monitor/utils"
	"senec-monitor/weather"
	"sync"
)

// Making this global is not ideal but passing this around in every function is worse imo
var LatestWeatherData *types.LatestWeather
var LatestTotal *types.LatestTotal
var LatestLocal *types.LatestLocal

func main() {
	logger := logging.NewLoggerWithoutFile()
	var (
		DbCreds      types.DbCreds
		UserCreds    types.UserInput
		WeatherCords types.Cordinate
		senec_ip     string
		err          error
	)
	if os.Getenv("mode") == "docker" {
		godotenv.Load("run/secrets/config")
		DbCreds, UserCreds, WeatherCords, senec_ip, err = utils.ReadEnvFromEnvFile()
	} else {
		godotenv.Load(".env")
		DbCreds, UserCreds, WeatherCords, senec_ip, err = utils.ReadEnvFromEnvFile()
	}
	if err != nil {
		logger.Fatal(err)
	}
	stmts, err := utils.ReadStatments()
	if err != nil {
		logger.Fatal(err)
	}
	service, err := db.NewPostgresService(DbCreds, stmts)
	if err != nil {
		logger.Fatal(err)
	}
	LatestWeatherData = &types.LatestWeather{Mu: sync.RWMutex{}, Data: types.ApiRespHourly{}}
	dataChan := make(chan *types.LocalApiDataWithCorrectTypes)

	LatestTotal = types.NewLatestTotal()
	go task.GetTotalEveryHour(UserCreds, logger, LatestTotal)

	LatestLocal = types.NewLatestLocal()
	weatherCh := make(chan types.ApiRespHourly)
	go task.CreateAndLoopLocalTask(logger, dataChan, senec_ip)
	go task.LoopAndUpdate(UserCreds, service, logger)

	go func() {
		for {
			select {
			case msg := <-weatherCh:
				{
					LatestWeatherData.Set(msg)
				}
			default:
				continue

			}
		}

	}()
	// get weather for each day in the morning
	go weather.GetWeatherEvery24(WeatherCords, logger, service)

	// for weather updates each hour
	go weather.GetWeatherHourly(WeatherCords, logger, weatherCh)
	server := server.NewServer(logger, service)

	server.Start(dataChan, LatestWeatherData, LatestTotal, LatestLocal)
}
