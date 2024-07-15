package utils

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"senec-monitor/types"
	"strconv"
	"strings"
)

func ReadStatments() (types.Statements, error) {
	typesForError := types.Statements{}
	statement1, err := readFile("statements/write.sql")
	if err != nil {
		return typesForError, err
	}
	statement2, err := readFile("statements/getData.sql")
	if err != nil {
		return typesForError, err
	}
	statement3, err := readFile("statements/getSpecificData.sql")
	if err != nil {
		return typesForError, err
	}
	statement4, err := readFile("statements/writeFromLocalApi.sql")
	if err != nil {
		return typesForError, err
	}
	statement5, err := readFile("statements/getDataFromLocal.sql")
	if err != nil {
		return typesForError, err
	}
	statement6, err := readFile("statements/writeDailyWeather.sql")
	if err != nil {
		return typesForError, err
	}
	res := types.Statements{
		Write:             statement1,
		GetData:           statement2,
		GetSpecificData:   statement3,
		WriteFromLocalApi: statement4,
		GetDataFromLocal:  statement5,
		WriteDailyWeather: statement6,
	}

	return res, nil

}
func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	line := ""
	i := 0
	for scanner.Scan() {
		if i == 0 {
			line = scanner.Text()
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err

	}
	return line, err

}
func ReadEnv(fileName string) (types.DbCreds, types.UserInput, types.Cordinate, string, error) {
	CordsForError := types.Cordinate{}
	CredsForError := types.DbCreds{}
	InputForError := types.UserInput{}
	file, err := os.Open(fileName)
	if err != nil {
		return CredsForError, InputForError, CordsForError, "", err
	}
	scanner := bufio.NewScanner(file)
	if err != nil {
		return CredsForError, InputForError, CordsForError, "", err
	}
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	CredsForReturn := &types.DbCreds{}
	InputForReturn := &types.UserInput{}
	CordsForReturn := &types.Cordinate{}
	configMap := map[string]string{
		"SENEC_PASS": "",
		"SENEC_USER": "",
		"DB_HOST":    "",
		"DB_PORT":    "",
		"DB_PASS":    "",
		"DB_USER":    "",
		"DB_NAME":    "",
		"LONGITUDE":  "",
		"LATITUDE":   "",
		"SENEC_IP":   "",
	}
	ip := ""

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		splits := strings.Split(line, "=")
		if len(splits) != 2 {
			return CredsForError, InputForError, CordsForError, "", fmt.Errorf("config file is wrong on line %v", i)
		}
		_, ok := configMap[splits[0]]
		if !ok {
			return CredsForError, InputForError, CordsForError, "", fmt.Errorf("%s is not a config option", splits[0])
		}
		configMap[splits[0]] = splits[1]
	}

	for k, v := range configMap {
		if v == "" {
			return CredsForError, InputForError, CordsForError, "", fmt.Errorf("missing value for %s", k)
		}
		if k == "DB_PASS" {
			CredsForReturn.Password = v
			continue
		}
		if k == "DB_PORT" {
			v_local, err := strconv.Atoi(v)
			if err != nil {
				return CredsForError, InputForError, CordsForError, "", fmt.Errorf("error parsing port into string")
			}
			CredsForReturn.Port = v_local

			continue
		}
		if k == "DB_USER" {
			CredsForReturn.User = v

			continue
		}
		if k == "DB_NAME" {
			CredsForReturn.DbName = v
			continue
		}
		if k == "DB_HOST" {
			CredsForReturn.Host = v
			continue
		}
		if k == "LONGITUDE" {
			CordsForReturn.Long = v
			continue
		}
		if k == "LATITUDE" {
			CordsForReturn.Lat = v
			continue
		}
		if k == "SENEC_PASS" {
			InputForReturn.Pass = v
			continue
		}
		if k == "SENEC_USER" {
			InputForReturn.User = v
			continue
		}
		if k == "SENEC_IP" {
			ip = v

		}

	}
	return *CredsForReturn, *InputForReturn, *CordsForReturn, ip, nil

}
func ReadEnvFromEnvFile() (types.DbCreds, types.UserInput, types.Cordinate, string, error) {
	CordsForError := types.Cordinate{}
	CredsForError := types.DbCreds{}
	InputForError := types.UserInput{}
	pass := os.Getenv("SENEC_PASS")
	user := os.Getenv("SENEC_USER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	db_pass := os.Getenv("DB_PASS")
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	long := os.Getenv("LONGITUDE")
	lat := os.Getenv("LATITUDE")
	ip := os.Getenv("SENEC_IP")
	if pass == "" || user == "" || host == "" || port == "" || db_pass == "" || db_name == "" || db_user == "" || long == "" || lat == "" || ip == "" {
		return CredsForError, InputForError, CordsForError, "", fmt.Errorf("missing value in env file")
	}
	port_int, err := strconv.Atoi(port)
	if err != nil {

		return CredsForError, InputForError, CordsForError, "", fmt.Errorf("error parsing int")
	}
	return types.DbCreds{
			Host:     host,
			User:     db_user,
			Password: db_pass,
			Port:     port_int,
			DbName:   db_name,
		}, types.UserInput{
			User: user,
			Pass: pass,
		}, types.Cordinate{
			Long: long,
			Lat:  lat,
		}, ip, nil

}
func ParseStringDataToStruct(data types.LocalApiResponse) (*types.LocalApiDataWithCorrectTypes, error) {
	stat_stat, err := parseHexStringU8(data.Energy.StatState)
	if err != nil {
		return nil, err
	}
	gui_bat_data_power, err := parseHexStringF32(data.Energy.GuiBatDataPower)
	if err != nil {
		return nil, err
	}
	gui_bat_data_power = gui_bat_data_power * 1e-3
	GUI_INVERTER_POWER, err := parseHexStringF32(data.Energy.GuiInverterPower)
	if err != nil {
		return nil, err
	}
	if GUI_INVERTER_POWER < 0 || GUI_INVERTER_POWER == -0 {
		GUI_INVERTER_POWER = 0
	}
	GUI_INVERTER_POWER = GUI_INVERTER_POWER * 1e-3
	GUI_HOUSE_POW, err := parseHexStringF32(data.Energy.GuiHousePow)
	if err != nil {
		return nil, err
	}
	GUI_HOUSE_POW = GUI_HOUSE_POW * 1e-3

	GUI_GRID_POW, err := parseHexStringF32(data.Energy.GuiGridPow)
	if err != nil {
		return nil, err
	}

	GUI_GRID_POW = GUI_GRID_POW * 1e-3

	GUI_BAT_DATA_FUEL_CHARGE, err := parseHexStringF32(data.Energy.GuiBatDataFuelCharge)
	if err != nil {
		return nil, err
	}

	GUI_CHARGING_INFO, err := parseHexStringU8(data.Energy.GuiChargingInfo)
	if err != nil {
		return nil, err
	}
	GUI_BOOSTING_INFO, err := parseHexStringU8(data.Energy.GuiBoostingInfo)
	if err != nil {

		return nil, err
	}

	return &types.LocalApiDataWithCorrectTypes{
		STAT_STATE:               stat_stat,
		GUI_BAT_DATA_POWER:       gui_bat_data_power,
		GUI_INVERTER_POWER:       GUI_INVERTER_POWER,
		GUI_HOUSE_POW:            GUI_HOUSE_POW,
		GUI_GRID_POW:             GUI_GRID_POW,
		GUI_BAT_DATA_FUEL_CHARGE: GUI_BAT_DATA_FUEL_CHARGE,
		GUI_CHARGING_INFO:        GUI_CHARGING_INFO,
		GUI_BOOSTING_INFO:        GUI_BOOSTING_INFO,
	}, nil
}
func parseHexStringU8(inp string) (uint8, error) {
	hexPart := strings.TrimLeft(inp, "u8_")

	n := new(big.Int)
	n.SetString(hexPart, 16)
	u_64 := n.Uint64()

	return uint8(u_64), nil

}
func parseHexStringF32(inp string) (float32, error) {
	hexPart := strings.TrimLeft(inp, "fl_")
	num, err := strconv.ParseUint(hexPart, 16, 32)
	if err != nil {
		return 0, err
	}
	floatVal := math.Float32frombits(uint32(num))

	return floatVal, nil

}
