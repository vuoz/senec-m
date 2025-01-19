package types

import (
	"fmt"
	pb "senec-monitor/proto"
	"strconv"
	"sync"
	"time"
)

type Cordinate struct {
	Long string `json:"long"`
	Lat  string `json:"lat"`
}

func (c *Cordinate) ToFloat() (Coordinates, error) {
	long, err := strconv.ParseFloat(c.Long, 64)
	if err != nil {
		return Coordinates{}, err
	}
	lat, err := strconv.ParseFloat(c.Lat, 64)
	if err != nil {
		return Coordinates{}, err
	}
	return Coordinates{Long: long, Lat: lat}, nil

}

type LatestWeather struct {
	Mu   sync.RWMutex
	Data ApiRespHourly
}
type LatestTotal struct {
	mu   sync.RWMutex
	data TotalData
	new  bool
}

type TotalData struct {
	Consumption string `json:"consumption"`
	Generated   string `json:"generated"`
	New         bool   `json:"new"`
}

func (d *TotalData) ToProto() *pb.TotalDataNew {
	return &pb.TotalDataNew{
		Consumption: d.Consumption,
		Generated:   d.Generated,
		New:         d.New,
	}

}

type LatestLocal struct {
	mu   sync.RWMutex
	data LocalApiDataWithCorrectTypes
}

func NewLatestLocal() *LatestLocal {
	return &LatestLocal{
		mu:   sync.RWMutex{},
		data: LocalApiDataWithCorrectTypes{},
	}
}
func (t *LatestLocal) Set(data LocalApiDataWithCorrectTypes) {
	t.mu.Lock()
	t.data = data
	t.mu.Unlock()
}
func (t *LatestLocal) Get() LocalApiDataWithCorrectTypes {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}

func NewLatestTotal() *LatestTotal {
	return &LatestTotal{
		mu:   sync.RWMutex{},
		data: TotalData{},
		new:  false,
	}

}

func (data *LatestTotal) Set(data_new TotalData) {
	data.mu.Lock()
	defer data.mu.Unlock()
	data.data = data_new
	data.new = true
	data.data.New = true

}
func (data *LatestTotal) Get() (TotalData, error) {
	data.mu.Lock()
	defer data.mu.Unlock()
	// since it only gets updated every hour
	if data.new {
		data.new = false
		copy_ := data.data
		data.data.New = false
		return copy_, nil
	} else {
		return data.data, fmt.Errorf("no new value")

	}

}

func (data *LatestWeather) Get() ApiRespHourly {
	data.Mu.RLock()
	defer data.Mu.RUnlock()
	return data.Data

}
func (data *LatestWeather) Set(toSet ApiRespHourly) {
	data.Mu.Lock()
	defer data.Mu.Unlock()
	data.Data = toSet
}

type DbCreds struct {
	Host     string
	Port     int
	User     string
	DbName   string
	Password string
}
type GetDataDB struct {
	Ts          int64   `json:"ts"`
	Generated   float64 `json:"generated"`
	Consumption float64 `json:"consumption"`
	Gridexport  float64 `json:"gridexport"`
	Accuexport  float64 `json:"accuexport"`
	Accuimport  float64 `json:"accuimport"`
	Gridimport  float64 `json:"gridimport"`
}
type UserInput struct {
	User string
	Pass string
}
type HandlerErrorResponse struct {
	Error string `json:"error"`
}
type Statements struct {
	Write             string
	GetData           string
	GetSpecificData   string
	WriteFromLocalApi string
	GetDataFromLocal  string
	WriteDailyWeather string
}
type StatusResponse struct {
	Wartungsplan struct {
		PossibleMaintenanceTypes []any `json:"possibleMaintenanceTypes"`
		MaintenanceDueSoon       bool  `json:"maintenanceDueSoon"`
		MaintenanceOverdue       bool  `json:"maintenanceOverdue"`
		MinorMaintenancePossible bool  `json:"minorMaintenancePossible"`
		Applicable               bool  `json:"applicable"`
	} `json:"wartungsplan"`
	SuppressedNotificationIds []any  `json:"suppressedNotificationIds"`
	SteuereinheitState        string `json:"steuereinheitState"`
	WartungNotwendig          bool   `json:"wartungNotwendig"`
	State                     int    `json:"state"`
	FirmwareVersion           int    `json:"firmwareVersion"`
	Lastupdated               int    `json:"lastupdated"`
	Powergenerated            struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"powergenerated"`
	Consumption struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"consumption"`
	Gridexport struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"gridexport"`
	Accuexport struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"accuexport"`
	Accuimport struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"accuimport"`
	Acculevel struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"acculevel"`
	McuOperationalModeID int `json:"mcuOperationalModeId"`
	Gridimport           struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"gridimport"`
	Machine string `json:"machine"`
}
type Response2 struct {
	Wartungsplan struct {
		PossibleMaintenanceTypes []any `json:"possibleMaintenanceTypes"`
		MaintenanceDueSoon       bool  `json:"maintenanceDueSoon"`
		MaintenanceOverdue       bool  `json:"maintenanceOverdue"`
		MinorMaintenancePossible bool  `json:"minorMaintenancePossible"`
		Applicable               bool  `json:"applicable"`
	} `json:"wartungsplan"`
	SuppressedNotificationIds []any  `json:"suppressedNotificationIds"`
	SteuereinheitState        string `json:"steuereinheitState"`
	WartungNotwendig          bool   `json:"wartungNotwendig"`
	State                     int    `json:"state"`
	FirmwareVersion           int    `json:"firmwareVersion"`
	Lastupdated               int    `json:"lastupdated"`
	Powergenerated            struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"powergenerated"`
	Consumption struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"consumption"`
	Gridexport struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"gridexport"`
	Accuexport struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"accuexport"`
	Accuimport struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"accuimport"`
	Acculevel struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"acculevel"`
	McuOperationalModeID int `json:"mcuOperationalModeId"`
	Gridimport           struct {
		Today float64 `json:"today"`
		Now   float64 `json:"now"`
	} `json:"gridimport"`
	Machine string `json:"machine"`
}
type Data24Hours struct {
	Val [][][]float64 `json:"val"`
}
type Accustatus struct {
	Akkutyp        string        `json:"akkutyp"`
	Fuelgauge      float64       `json:"fuelgauge"`
	Batterycurrent float64       `json:"batterycurrent"`
	Batteryvoltage float64       `json:"batteryvoltage"`
	Bars           string        `json:"bars"`
	Lastupdated    int64         `json:"lastupdated"`
	Val            [][][]float64 `json:"val"`
	Capacity       int           `json:"capacity"`
}

type LocalApiResponse struct {
	Energy struct {
		StatState            string `json:"STAT_STATE"`
		GuiBatDataPower      string `json:"GUI_BAT_DATA_POWER"`
		GuiInverterPower     string `json:"GUI_INVERTER_POWER"`
		GuiHousePow          string `json:"GUI_HOUSE_POW"`
		GuiGridPow           string `json:"GUI_GRID_POW"`
		GuiBatDataFuelCharge string `json:"GUI_BAT_DATA_FUEL_CHARGE"`
		GuiChargingInfo      string `json:"GUI_CHARGING_INFO"`
		GuiBoostingInfo      string `json:"GUI_BOOSTING_INFO"`
	} `json:"ENERGY"`
	SysUpdate struct {
		UpdateAvailable string `json:"UPDATE_AVAILABLE"`
	} `json:"SYS_UPDATE"`
	Steca struct {
		Au2020VersionMismatch string   `json:"AU2020_VERSION_MISMATCH"`
		Bat                   string   `json:"BAT"`
		BdcState              []string `json:"BDC_STATE"`
		Error                 string   `json:"ERROR"`
		Errortext             string   `json:"ERRORTEXT"`
		Island                string   `json:"ISLAND"`
		NumPvConfigPossible   string   `json:"NUM_PV_CONFIG_POSSIBLE"`
		Pv                    string   `json:"PV"`
		Pvss                  string   `json:"PVSS"`
		PvConfigPossible      []string `json:"PV_CONFIG_POSSIBLE"`
		PvInputs              string   `json:"PV_INPUTS"`
		Relays                string   `json:"RELAYS"`
		Startup               string   `json:"STARTUP"`
		StartupAdd            string   `json:"STARTUP_ADD"`
	} `json:"STECA"`
}

type LocalApiDataWithCorrectTypes struct {
	STAT_STATE               uint8
	GUI_BAT_DATA_POWER       float32
	GUI_INVERTER_POWER       float32
	GUI_HOUSE_POW            float32
	GUI_GRID_POW             float32
	GUI_BAT_DATA_FUEL_CHARGE float32
	GUI_CHARGING_INFO        uint8
	GUI_BOOSTING_INFO        uint8
}

type LocalApiDataWithCorrectTypesWithTimeStamp struct {
	TS                 int64   `json:"ts"`
	STAT_STATE         uint8   `json:"stat_state"`
	GUI_BAT_DATA_POWER float32 `json:"gui_bat_data_power"`
	GUI_INVERTER_POWER float32 `json:"gui_inverter_power"`
	GUI_HOUSE_POW      float32 `json:"gui_house_pow"`

	// this is the grid_export & import in one
	GUI_GRID_POW             float32 `json:"gui_grid_pow"`
	GUI_BAT_DATA_FUEL_CHARGE float32 `json:"gui_bat_data_fuel_charge"`
	GUI_CHARGING_INFO        uint8   `json:"gui_charging_info"`
	GUI_BOOSTING_INFO        uint8   `json:"gui_boosting_info"`
}

// crazy type name ik
type LocalApiDataWithCorrectTypesWithTimeStampStringsWithWeather struct {
	TS                       string               `json:"ts"`
	STAT_STATE               string               `json:"stat_state"`
	GUI_BAT_DATA_POWER       string               `json:"gui_bat_data_power"`
	GUI_INVERTER_POWER       string               `json:"gui_inverter_power"`
	GUI_HOUSE_POW            string               `json:"gui_house_pow"`
	GUI_GRID_POW             string               `json:"gui_grid_pow"`
	GUI_BAT_DATA_FUEL_CHARGE string               `json:"gui_bat_data_fuel_charge"`
	GUI_CHARGING_INFO        string               `json:"gui_charging_info"`
	GUI_BOOSTING_INFO        string               `json:"gui_boosting_info"`
	Weather                  ApiRespHourlyStrings `json:"weather"`
	Total_data               TotalData            `json:"total_data"`
}

func (data *LocalApiDataWithCorrectTypes) ConvertToStrings(weather ApiRespHourly, total_data TotalData) LocalApiDataWithCorrectTypesWithTimeStampStringsWithWeather {
	ts := time.Now().Unix()
	return LocalApiDataWithCorrectTypesWithTimeStampStringsWithWeather{
		TS:                       time.Unix(ts, 0).Format(time.Kitchen),
		STAT_STATE:               strconv.FormatUint(uint64(data.STAT_STATE), 10),
		GUI_BAT_DATA_POWER:       strconv.FormatFloat(float64(data.GUI_BAT_DATA_POWER), 'f', 2, 32),
		GUI_INVERTER_POWER:       strconv.FormatFloat(float64(data.GUI_INVERTER_POWER), 'f', 2, 32),
		GUI_HOUSE_POW:            strconv.FormatFloat(float64(data.GUI_HOUSE_POW), 'f', 2, 32),
		GUI_GRID_POW:             strconv.FormatFloat(float64(data.GUI_GRID_POW), 'f', 2, 32),
		GUI_BAT_DATA_FUEL_CHARGE: strconv.FormatFloat(float64(data.GUI_BAT_DATA_FUEL_CHARGE), 'f', 1, 32),
		GUI_CHARGING_INFO:        strconv.FormatUint(uint64(data.GUI_CHARGING_INFO), 10),
		GUI_BOOSTING_INFO:        strconv.FormatUint(uint64(data.GUI_BOOSTING_INFO), 10),
		Total_data:               total_data,
		// will need to do some preprocessing before it will be send down to client since it is too much data and alot of it isn't needed
		Weather: weather.ToStructOfStrings(),
	}

}
func (data *LocalApiDataWithCorrectTypesWithTimeStampStringsWithWeather) ConvertToProto() pb.NewUiStruct {
	return pb.NewUiStruct{
		//generate all the binding
		Ts:                   data.TS,
		StatState:            data.STAT_STATE,
		GuiBatDataPower:      data.GUI_BAT_DATA_POWER,
		GuiInverterPower:     data.GUI_INVERTER_POWER,
		GuiHousePow:          data.GUI_HOUSE_POW,
		GuiGridPow:           data.GUI_GRID_POW,
		GuiBatDataFuelCharge: data.GUI_BAT_DATA_FUEL_CHARGE,
		GuiChargingInfo:      data.GUI_CHARGING_INFO,
		GuiBoostingInfo:      data.GUI_BOOSTING_INFO,
		Weather:              data.Weather.ConvertToProto(),
		TotalData:            data.Total_data.ToProto(),
	}

}

type WeatherApiResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	HourlyUnits          struct {
		Time       string `json:"time"`
		CloudCover string `json:"cloud_cover"`
		Visibility string `json:"visibility"`
	} `json:"hourly_units"`
	Hourly struct {
		Time       []string  `json:"time"`
		CloudCover []int     `json:"cloud_cover"`
		Visibility []float64 `json:"visibility"`
	} `json:"hourly"`
	DailyUnits struct {
		Time               string `json:"time"`
		Sunrise            string `json:"sunrise"`
		Sunset             string `json:"sunset"`
		SunshineDuration   string `json:"sunshine_duration"`
		UvIndexMax         string `json:"uv_index_max"`
		UvIndexClearSkyMax string `json:"uv_index_clear_sky_max"`
	} `json:"daily_units"`
	Daily struct {
		Time               []string  `json:"time"`
		Sunrise            []string  `json:"sunrise"`
		Sunset             []string  `json:"sunset"`
		SunshineDuration   []float64 `json:"sunshine_duration"`
		UvIndexMax         []float64 `json:"uv_index_max"`
		UvIndexClearSkyMax []float64 `json:"uv_index_clear_sky_max"`
	} `json:"daily"`
}

type WeatherDailyResp struct {
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	GenerationtimeMS     float64    `json:"generationtime_ms"`
	UTCOffsetSeconds     float64    `json:"utc_offset_seconds"`
	Timezone             string     `json:"timezone"`
	TimezoneAbbreviation string     `json:"timezone_abbreviation"`
	Elevation            float64    `json:"elevation"`
	DailyUnits           DailyUnits `json:"daily_units"`
	Daily                Daily      `json:"daily"`
}

func (data *WeatherDailyResp) ToDbType() WeatherDailyDb {
	date := time.Now().Format(time.DateOnly)
	return WeatherDailyDb{
		Date:  date,
		Daily: data.Daily,
	}
}

type WeatherDailyDb struct {
	Date  string
	Daily Daily
}

type Daily struct {
	Time               []string  `json:"time"`
	WeatherCode        []float64 `json:"weather_code"`
	Temperature2MMax   []float64 `json:"temperature_2m_max"`
	Temperature2MMin   []float64 `json:"temperature_2m_min"`
	Sunrise            []string  `json:"sunrise"`
	Sunset             []string  `json:"sunset"`
	DaylightDuration   []float64 `json:"daylight_duration"`
	SunshineDuration   []float64 `json:"sunshine_duration"`
	UvIndexMax         []float64 `json:"uv_index_max"`
	UvIndexClearSkyMax []float64 `json:"uv_index_clear_sky_max"`
}

type DailyUnits struct {
	Time               string `json:"time"`
	WeatherCode        string `json:"weather_code"`
	Temperature2MMax   string `json:"temperature_2m_max"`
	Temperature2MMin   string `json:"temperature_2m_min"`
	Sunrise            string `json:"sunrise"`
	Sunset             string `json:"sunset"`
	DaylightDuration   string `json:"daylight_duration"`
	SunshineDuration   string `json:"sunshine_duration"`
	UvIndexMax         string `json:"uv_index_max"`
	UvIndexClearSkyMax string `json:"uv_index_clear_sky_max"`
}
type ApiRespHourlyStrings struct {
	Hourly HourlyForRespHourlyStrings `json:"hourly"`
	Daily  DailyHourly                `json:"daily"`
}

func (d *ApiRespHourlyStrings) ConvertToProto() *pb.WeatherNew {
	return &pb.WeatherNew{
		Hourly: d.Hourly.ConvertToProto(),
		Daily:  d.Daily.ConvertToProto(),
	}
}

type HourlyForRespHourlyStrings struct {
	Time            []string `json:"time"`
	Temperature2M   []string `json:"temperature_2m"`
	Rain            []string `json:"rain"`
	Showers         []string `json:"showers"`
	CloudCover      []string `json:"cloud_cover"`
	UvIndex         []string `json:"uv_index"`
	UvIndexClearSky []string `json:"uv_index_clear_sky"`
}

func (d *HourlyForRespHourlyStrings) ConvertToProto() *pb.HourlyNew {
	return &pb.HourlyNew{
		Time:            d.Time,
		Temperature_2M:  d.Temperature2M,
		Rain:            d.Rain,
		Showers:         d.Showers,
		CloudCover:      d.CloudCover,
		UvIndex:         d.UvIndex,
		UvIndexClearSky: d.UvIndexClearSky,
	}
}

type ApiRespHourly struct {
	Latitude             float64                  `json:"latitude"`
	Longitude            float64                  `json:"longitude"`
	GenerationtimeMS     float64                  `json:"generationtime_ms"`
	UTCOffsetSeconds     float64                  `json:"utc_offset_seconds"`
	Timezone             string                   `json:"timezone"`
	TimezoneAbbreviation string                   `json:"timezone_abbreviation"`
	Elevation            float64                  `json:"elevation"`
	HourlyUnits          HourlyUnitsForRespHourly `json:"hourly_units"`
	Hourly               HourlyForRespHourly      `json:"hourly"`
	DailyUnits           DailyUnitsHourly         `json:"daily_units"`
	Daily                DailyHourly              `json:"daily"`
}

func (data *ApiRespHourly) ToStructOfStrings() ApiRespHourlyStrings {
	data.Daily.StripDate()
	return ApiRespHourlyStrings{
		Hourly: data.Hourly.ToStructOfStrings(),
		Daily:  data.Daily,
	}

}

type DailyHourly struct {
	Time    []string `json:"time"`
	Sunset  []string `json:"sunset"`
	Sunrise []string `json:"sunrise"`
}

func (d *DailyHourly) ConvertToProto() *pb.DailyNew {
	return &pb.DailyNew{
		Time:    d.Time,
		Sunset:  d.Sunset,
		Sunrise: d.Sunrise,
	}
}
func (data *DailyHourly) StripDate() {
	new_sunrise := make([]string, len(data.Sunrise))
	for i, val := range data.Sunrise {
		if len(val) < 11 {
			continue
		}
		new_sunrise[i] = val[11:]
	}
	new_sunset := make([]string, len(data.Sunset))
	for i, val := range data.Sunset {
		if len(val) < 11 {
			continue
		}
		new_sunset[i] = val[11:]
	}
	data.Sunrise = new_sunrise
	data.Sunset = new_sunset

}

type DailyUnitsHourly struct {
	Time    string `json:"time"`
	Sunset  string `json:"sunset"`
	Sunrise string `json:"sunrise"`
}
type HourlyForRespHourly struct {
	Time            []string  `json:"time"`
	Temperature2M   []float64 `json:"temperature_2m"`
	Rain            []float64 `json:"rain"`
	Showers         []float64 `json:"showers"`
	CloudCover      []float64 `json:"cloud_cover"`
	UvIndex         []float64 `json:"uv_index"`
	UvIndexClearSky []float64 `json:"uv_index_clear_sky"`
}

func (data *HourlyForRespHourly) ToStructOfStrings() HourlyForRespHourlyStrings {
	temp2M := make([]string, len(data.Temperature2M))
	for i, val := range data.Temperature2M {
		temp2M[i] = strconv.FormatFloat(val, 'f', 1, 64)
	}

	rain := make([]string, len(data.Rain))
	for i, val := range data.Rain {
		rain[i] = strconv.FormatFloat(val, 'f', 1, 64)
	}

	showers := make([]string, len(data.Showers))
	for i, val := range data.Showers {
		showers[i] = strconv.FormatFloat(val, 'f', 1, 64)
	}

	cloudCover := make([]string, len(data.CloudCover))
	for i, val := range data.CloudCover {
		cloudCover[i] = strconv.FormatFloat(val, 'f', 1, 64)
	}

	uvIndex := make([]string, len(data.UvIndex))
	for i, val := range data.UvIndex {
		uvIndex[i] = strconv.FormatFloat(val, 'f', 1, 64)
	}

	uvIndexClearSky := make([]string, len(data.UvIndexClearSky))
	for i, val := range data.UvIndexClearSky {
		uvIndexClearSky[i] = strconv.FormatFloat(val, 'f', 1, 64)
	}

	return HourlyForRespHourlyStrings{
		Time:            data.Time,
		Temperature2M:   temp2M,
		Rain:            rain,
		Showers:         showers,
		CloudCover:      cloudCover,
		UvIndex:         uvIndex,
		UvIndexClearSky: uvIndexClearSky,
	}
}

type HourlyUnitsForRespHourly struct {
	Time            string `json:"time"`
	Temperature2M   string `json:"temperature_2m"`
	Rain            string `json:"rain"`
	Showers         string `json:"showers"`
	CloudCover      string `json:"cloud_cover"`
	UvIndex         string `json:"uv_index"`
	UvIndexClearSky string `json:"uv_index_clear_sky"`
}
type PredictionResponse struct {
	Data []float64 `json:"data"`
}

type PredictionRequest struct {
	Date  string      `json:"date"`
	Coord Coordinates `json:"coordinates"`
}
type Coordinates struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}
