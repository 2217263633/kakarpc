package cusrequest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PolicyType int32

const (
	Get    PolicyType = 0
	Post   PolicyType = 1
	Put    PolicyType = 2
	DELETE PolicyType = 3
)

var methods = []string{"GET", "POST", "PUT", "DELETE"}

func Request(_url string, method PolicyType, data any, Authorization string) (map[string]interface{}, error) {
	// url := "https://data-sk-wh.scasst.net/whsj/pc/account"
	// method := "POST"
	var resu map[string]interface{} = make(map[string]interface{})

	jsonVal, err := json.Marshal(data)

	payload := strings.NewReader(string(jsonVal))
	if err != nil {
		payload = nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(methods[method], _url, payload)
	if err != nil {
		return resu, err
	}
	req.Header.Add("Authorization", Authorization)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	if method != Get {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return resu, err
	}

	if res.StatusCode != 200 {
		return resu, fmt.Errorf(fmt.Sprintf("status code is %d", res.StatusCode))
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		// logger.Info(err, "33333")
		return resu, err
	}

	err = json.Unmarshal(body, &resu)
	return resu, err
}

func RequestWithHeader(_url string, method PolicyType, data any, headers http.Header) (map[string]interface{}, error) {
	var resu map[string]interface{} = make(map[string]interface{})
	jsonVal, err := json.Marshal(data)
	payload := strings.NewReader(string(jsonVal))
	if err != nil {
		payload = nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(methods[method], _url, payload)
	if err != nil {
		return resu, err
	}
	req.Header = headers
	res, err := client.Do(req)
	if err != nil {
		return resu, err
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("status code is %d", res.StatusCode)
		return resu, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return resu, err
	}
	err = json.Unmarshal(body, &resu)
	return resu, err
}

// return string error
func Request2(_url string, method PolicyType, data any, Authorization string) (string, error) {
	var resu string = ""
	jsonVal, err := json.Marshal(data)
	payload := strings.NewReader(string(jsonVal))

	if err != nil {
		payload = nil
	}

	client := &http.Client{
		Timeout: time.Second * 25,
	}
	req, err := http.NewRequest(methods[method], _url, payload)
	if err != nil {

		return resu, err
	}
	req.Header.Add("Authorization", Authorization)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	if method != Get {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		return resu, err
	}

	if res.StatusCode != 200 {
		return resu, fmt.Errorf(fmt.Sprintf("status code is %d", res.StatusCode))
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return resu, err
	}
	return string(body), nil
}
