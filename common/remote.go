package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	GeneralRemoteCallTimeout = 10 // seconds
)

var logger = log.GetLogger()

//=============================================================
//
//=============================================================

var httpClient = &http.Client{
	Timeout: time.Duration(GeneralRemoteCallTimeout) * time.Second,
}

func RemoteCallWithBody(method, url string, token, user string, body []byte, contentType string) (*http.Response, []byte, error) {
	logger.Info("method: %s, url: %s, token: %s, contentType: %s, body: %s", method, url, token, contentType, string(body))

	var request *http.Request
	var err error
	if len(body) == 0 {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, bytes.NewReader(body))
	}
	if err != nil {
		return nil, nil, err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if token != "" {
		request.Header.Set("Authorization", token)
	}
	if user != "" {
		request.Header.Set("User", user)
	}
	//client := &http.Client{
	//	Timeout: time.Duration(GeneralRemoteCallTimeout) * time.Second,
	//}
	response, err := httpClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	return response, bytes, err
}

func RemoteCallWithJsonBody(method string, url string, token, user string, jsonBody []byte) (*http.Response, []byte, error) {
	return RemoteCallWithBody(method, url, token, user, jsonBody, "application/json; charset=utf-8")
}

func RemoteCall(method string, url string, token, user string) (*http.Response, []byte, error) {
	return RemoteCallWithBody(method, url, token, user, nil, "")
}

func GetRequestData(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseRequestJsonAsMap(r *http.Request) (map[string]interface{}, error) {
	data, err := GetRequestData(r)
	if err != nil {
		return nil, err
	}

	m, err := ParseJsonToMap(data)
	if err != nil {
		logger.Error("ParseJsonToMap r.Body (%s) error: %s", string(data), err.Error())
	}

	return m, err
}

func ParseRequestJsonInto(r *http.Request, into interface{}) error {
	data, err := GetRequestData(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, into)
}

func ParseRequestJsonIntoWithValidateParams(r *http.Request, correctInput []string, into interface{}) error {
	data, err := GetRequestData(r)
	if err != nil {
		logger.Error("Get request data err: %v", err)
		return err
	}

	paramsMap := make(map[string]interface{})
	err = json.Unmarshal(data, &paramsMap)
	if err != nil {
		logger.Error("Unmarshal err: %v", err)
		return err
	}

	for _, value := range correctInput {
		if paramsMap[value] == nil {
			return errors.New("input params not correct")
		}
	}

	return json.Unmarshal(data, into)
}
