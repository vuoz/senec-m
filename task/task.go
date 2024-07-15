package task

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"senec-monitor/db"
	"senec-monitor/logging"
	"senec-monitor/scheduler"
	"senec-monitor/types"
	"strconv"
	"strings"
	"time"
)

const (
	ModeOverview = iota
	Mode24Hours  = iota
)

type Task struct {
	c   *http.Client
	log logging.Logger
}

func (t *Task) getData() (*types.Response2, error) {
	headers := http.Header{
		"Accept":                    []string{"application/json, text/plain, */*"},
		"Accept-Encoding":           []string{"gzip, deflate, br"},
		"Accept-Language":           []string{"de-DE,de;q=0.9"},
		"Cache-Control":             []string{"no-cache"},
		"Connection":                []string{"keep-alive"},
		"Host":                      []string{"mein-senec.de"},
		"Origin":                    []string{"https://mein-senec.de"},
		"Pragma":                    []string{"no-cache"},
		"Referer":                   []string{"https://mein-senec.de/endkunde"},
		"Sec-Ch-Ua":                 []string{`"Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"`},
		"Sec-Ch-Ua-Mobile":          []string{"?0"},
		"Sec-Ch-Ua-Platform":        []string{"Windows"},
		"Sec-Fetch-Dest":            []string{"document"},
		"Sec-Fetch-Mode":            []string{"navigate"},
		"Sec-Fetch-Site":            []string{"same-origin"},
		"Sec-Fetch-User":            []string{"?1"},
		"Upgrade-Insecure-Requests": []string{"1"},
		"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
	}

	req, err := http.NewRequest("GET", "https://mein-senec.de/endkunde/api/status/getstatusoverview.php?anlageNummer=0", nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	rsp, err := t.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("status error %d", rsp.StatusCode)
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if strings.Index(string(body), "<") == 0 {
		return nil, fmt.Errorf("login session expired")
	}

	var data *types.Response2
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil

}
func (t *Task) sendLogin(email string, pass string) error {
	headers := http.Header{
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Encoding":           []string{"gzip, deflate, br"},
		"Accept-Language":           []string{"de-DE,de;q=0.9"},
		"Cache-Control":             []string{"no-cache"},
		"Connection":                []string{"keep-alive"},
		"Content-Type":              []string{"application/x-www-form-urlencoded"},
		"Host":                      []string{"mein-senec.de"},
		"Origin":                    []string{"https://mein-senec.de"},
		"Pragma":                    []string{"no-cache"},
		"Referer":                   []string{"https://mein-senec.de/auth/login"},
		"Sec-Ch-Ua":                 []string{`"Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"`},
		"Sec-Ch-Ua-Mobile":          []string{"?0"},
		"Sec-Ch-Ua-Platform":        []string{"Windows"},
		"Sec-Fetch-Dest":            []string{"document"},
		"Sec-Fetch-Mode":            []string{"navigate"},
		"Sec-Fetch-Site":            []string{"same-origin"},
		"Sec-Fetch-User":            []string{"?1"},
		"Upgrade-Insecure-Requests": []string{"1"},
		"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
	}
	data := url.Values{}
	data.Set("username", email)
	data.Set("password", pass)
	reqBody := data.Encode()

	req, err := http.NewRequest("POST", "https://mein-senec.de/auth/login", strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header = headers
	resp, err := t.c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error %d", resp.StatusCode)
	}

	return nil

}
func (t *Task) getLogin() error {

	headers := http.Header{
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Encoding":           []string{"gzip, deflate, br"},
		"Accept-Language":           []string{"de-DE,de;q=0.9"},
		"Cache-Control":             []string{"no-cache"},
		"Connection":                []string{"keep-alive"},
		"Host":                      []string{"mein-senec.de"},
		"Pragma":                    []string{"no-cache"},
		"Sec-Ch-Ua":                 []string{`"Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"`},
		"Sec-Ch-Ua-Mobile":          []string{"?0"},
		"Sec-Ch-Ua-Platform":        []string{"Windows"},
		"Sec-Fetch-Dest":            []string{"document"},
		"Sec-Fetch-Mode":            []string{"navigate"},
		"Sec-Fetch-Site":            []string{"same-origin"},
		"Sec-Fetch-User":            []string{"?1"},
		"Upgrade-Insecure-Requests": []string{"1"},
		"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
	}
	req, err := http.NewRequest("GET", "https://mein-senec.de/auth/login", nil)
	if err != nil {
		return err
	}
	req.Header = headers
	resp, err := t.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status error %d", resp.StatusCode)
	}
	return nil

}

func (t *Task) run(inpt types.UserInput) *types.Response2 {

	for {
		err := t.getLogin()
		if err != nil {
			t.log.Err("Eroror getting login page")
			time.Sleep(5 * time.Second)
			continue
		}
		t.log.Success("Got login page")
		break

	}
	for {
		err := t.sendLogin(inpt.User, inpt.Pass)
		if err != nil {
			t.log.Err("Error sending login")
			time.Sleep(5 * time.Second)
			continue
		}
		t.log.Success("Succesfully send login")
		break
	}

	var res *types.Response2
	for {
		var err error
		res, err = t.getData()
		if err != nil {
			t.log.Err("Error getting data ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		t.log.Success("Got Data")
		break
	}
	return res

}
func LaunchTask(inpt types.UserInput, log logging.Logger) (*types.Response2, error) {
	checkredirect := func(req *http.Request, via []*http.Request) error {
		redirects := 0
		return func(req *http.Request, via []*http.Request) error {
			if redirects == 12 {
				return fmt.Errorf("to many redirects: %d", redirects)
			}
			redirects++
			return nil

		}(req, via)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Jar:           jar,
		CheckRedirect: checkredirect,
	}
	task := &Task{
		c:   &client,
		log: log,
	}

	result := task.run(inpt)
	return result, nil

}
func newTask(log logging.Logger) (*Task, error) {
	checkredirect := func(req *http.Request, via []*http.Request) error {
		redirects := 0
		return func(req *http.Request, via []*http.Request) error {
			if redirects == 12 {
				return fmt.Errorf("to many redirects: %d", redirects)
			}
			redirects++
			return nil

		}(req, via)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Jar:           jar,
		CheckRedirect: checkredirect,
	}
	task := &Task{
		c:   &client,
		log: log,
	}
	return task, nil

}
func LoopAndUpdate(UserCreds types.UserInput, service db.DbService, logger logging.Logger) {

	for {
		timeToSleep := scheduler.Scheduler()
		time.Sleep(timeToSleep)
		result, err := LaunchTask(UserCreds, logger)
		if err != nil {
			logger.Fatal(err)
		}

		if err := service.Write(result); err != nil {
			logger.Fatal(err)
		}
		logger.Success("updated database")
		// To ensure that the task sleeps until after 00:00 and doesnt enter a loop
		time.Sleep(10 * time.Minute)

	}

}

func (t *Task) getTotalEveryHour(creds types.UserInput, logger logging.Logger, latestTotal *types.LatestTotal) error {
	task := t
	retries := 0
	for {
		err := task.getLogin()
		if err != nil {
			logger.Err("Error logging in task for total data err: ", err)
			retries++
			if retries > 5 {
				logger.Err("Hit 6 retries for total task trying again later!")
				return err
			}
			continue
		}
		break

	}
	retries = 0
	for {
		err := task.sendLogin(creds.User, creds.Pass)
		if err != nil {
			if retries > 5 {
				logger.Err("Hit 6 retries for total task trying again later!")
				return err
			}
			logger.Err("Error logging in task for total data err: ", err)
			retries++
			continue

		}
		break

	}
	retries = 0
	for {
		data, err := task.getData()
		if err != nil {
			retries++
			if retries > 5 {
				logger.Err("Hit retry limit for total task trying again later err: ", err)
				return err
			}

			logger.Err("Error data for total data task retry: ", retries, " err: ", err)
			continue
		}
		retries = 0
		latestTotal.Set(types.TotalData{
			Consumption: strconv.FormatFloat(float64(data.Consumption.Today), 'f', 2, 32),
			Generated:   strconv.FormatFloat(float64(data.Powergenerated.Today), 'f', 2, 32),
			New:         true,
		})
		break

	}
	return nil

}
func GetTotalEveryHour(creds types.UserInput, logger logging.Logger, latestTotal *types.LatestTotal) {
	task, err := newTask(logger)
	if err != nil {
		logger.Fatal("Error creating task")
	}
	totalRetries := 0
	for {
		if err := task.getTotalEveryHour(creds, logger, latestTotal); err != nil {
			totalRetries++
			if totalRetries >= 2 {
				logger.Err(" Error hit limit for local task")
			}
			time.Sleep(10 * time.Minute)
			continue
		}
		time.Sleep(1 * time.Hour)
		continue

	}

}
