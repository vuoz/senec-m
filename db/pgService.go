package db

import (
	"reflect"
	"time"

	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"senec-monitor/types"
)

type PostgresService struct {
	db         *sql.DB
	statements types.Statements
}

func NewPostgresService(creds types.DbCreds, stmts types.Statements) (DbService, error) {
	connectstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", creds.Host, creds.Port, creds.User, creds.Password, creds.DbName)
	db, err := sql.Open("postgres", connectstr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresService{
		db:         db,
		statements: stmts,
	}, nil

}
func (pg *PostgresService) WriteWeather24Hours(data types.WeatherDailyDb) error {
	v := reflect.ValueOf(data.Daily)

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		if fieldVal.Kind() != reflect.Slice || fieldVal.Kind() == reflect.Array {
			return fmt.Errorf("value not correct for field: %s", fieldVal.String())
		}
		if fieldVal.Len() < 1 {
			return fmt.Errorf("no value present for field: %s", fieldVal.String())
		}

	}
	daily := data.Daily

	res, err := pg.db.Query(pg.statements.WriteDailyWeather, data.Date, daily.WeatherCode[0], daily.Temperature2MMax[0], daily.Temperature2MMin[0], daily.Sunrise[0], daily.Sunset[0], daily.DaylightDuration[0], daily.SunshineDuration[0], daily.UvIndexMax[0], daily.UvIndexClearSkyMax[0])
	if err != nil {
		return err
	}
	defer res.Close()

	return nil
}
func (pg *PostgresService) Write(data *types.Response2) error {
	ts := time.Now().Unix()
	res, err := pg.db.Query(pg.statements.Write, ts, data.Powergenerated.Today, data.Consumption.Today, data.Gridexport.Today, data.Accuexport.Today, data.Accuimport.Today, data.Gridimport.Today)
	if err != nil {
		return err
	}
	res.Close()
	return nil

}
func (pg *PostgresService) GetData() ([]types.GetDataDB, error) {
	typesForError := []types.GetDataDB{}
	res, err := pg.db.Query(pg.statements.GetData)
	if err != nil {
		return typesForError, nil
	}
	data := []types.GetDataDB{}
	for res.Next() {

		var (
			Ts          int64
			Generated   float64
			Consumption float64
			Gridexport  float64
			Accuexport  float64
			Accuimport  float64
			Gridimport  float64
		)
		if err := res.Scan(&Ts, &Generated, &Consumption, &Gridexport, &Accuexport, &Accuimport, &Gridimport); err != nil {
			return typesForError, err
		}
		newData := types.GetDataDB{
			Ts:          Ts,
			Generated:   Generated,
			Consumption: Consumption,
			Gridexport:  Gridexport,
			Accuexport:  Accuexport,
			Accuimport:  Accuimport,
			Gridimport:  Gridimport,
		}
		data = append(data, newData)

	}
	res.Close()

	return data, nil

}
func (pg *PostgresService) GetSpecificData(ts int64) (types.GetDataDB, error) {
	typesForError := types.GetDataDB{}
	res, err := pg.db.Query(pg.statements.GetSpecificData, ts)
	if err != nil {
		return typesForError, err
	}
	results := []types.GetDataDB{}
	for res.Next() {

		var (
			Ts          int64
			Generated   float64
			Consumption float64
			Gridexport  float64
			Accuexport  float64
			Accuimport  float64
			Gridimport  float64
		)
		if err := res.Scan(&Ts, &Generated, &Consumption, &Gridexport, &Accuexport, &Accuimport, &Gridimport); err != nil {
			return typesForError, err
		}
		newData := types.GetDataDB{
			Ts:          Ts,
			Generated:   Generated,
			Consumption: Consumption,
			Gridexport:  Gridexport,
			Accuexport:  Accuexport,
			Accuimport:  Accuimport,
			Gridimport:  Gridimport,
		}
		results = append(results, newData)
	}

	res.Close()
	if len(results) == 0 {
		return typesForError, NewUserInputError("there are no results for this timestamp")
	}
	if len(results) > 1 {
		return typesForError, NewUserInputError("there is more than one record for the timestamp")
	}
	return results[0], nil

}
func (pg *PostgresService) WriteLocalApiData(data types.LocalApiDataWithCorrectTypes) error {
	ts := time.Now().Unix()
	res, err := pg.db.Query(pg.statements.WriteFromLocalApi, ts, data.STAT_STATE, data.GUI_BAT_DATA_POWER, data.GUI_INVERTER_POWER, data.GUI_HOUSE_POW, data.GUI_GRID_POW, data.GUI_BAT_DATA_FUEL_CHARGE, data.GUI_CHARGING_INFO, data.GUI_BOOSTING_INFO)
	if err != nil {
		return err
	}

	res.Close()
	return nil
}
func (pg *PostgresService) GetLatestFromLocal() (types.LocalApiDataWithCorrectTypesWithTimeStamp, error) {
	typesForError := types.LocalApiDataWithCorrectTypesWithTimeStamp{}
	res, err := pg.db.Query(pg.statements.GetDataFromLocal)
	if err != nil {
		return typesForError, err
	}

	results := []types.LocalApiDataWithCorrectTypesWithTimeStamp{}
	for res.Next() {
		var (
			STAT_STATE               uint8
			GUI_BAT_DATA_POWER       float32
			GUI_INVERTER_POWER       float32
			GUI_HOUSE_POW            float32
			GUI_GRID_POW             float32
			GUI_BAT_DATA_FUEL_CHARGE float32
			GUI_CHARGING_INFO        uint8
			GUI_BOOSTING_INFO        uint8
			TS                       int64
		)
		if err := res.Scan(&STAT_STATE, &GUI_BAT_DATA_POWER, &GUI_INVERTER_POWER, &GUI_HOUSE_POW, &GUI_GRID_POW, &GUI_BAT_DATA_FUEL_CHARGE, &GUI_CHARGING_INFO, &GUI_BOOSTING_INFO, &TS); err != nil {
			return typesForError, err
		}
		newData := types.LocalApiDataWithCorrectTypesWithTimeStamp{

			STAT_STATE:               STAT_STATE,
			GUI_BAT_DATA_POWER:       GUI_BAT_DATA_POWER,
			GUI_INVERTER_POWER:       GUI_INVERTER_POWER,
			GUI_HOUSE_POW:            GUI_HOUSE_POW,
			GUI_GRID_POW:             GUI_GRID_POW,
			GUI_BAT_DATA_FUEL_CHARGE: GUI_BAT_DATA_FUEL_CHARGE,
			GUI_CHARGING_INFO:        GUI_CHARGING_INFO,
			GUI_BOOSTING_INFO:        GUI_BOOSTING_INFO,
			TS:                       TS,
		}
		results = append(results, newData)

	}

	res.Close()
	if len(results) == 0 {
		return typesForError, NewUserInputError("no data available")
	}
	return results[0], nil

}
