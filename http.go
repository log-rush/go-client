package logRushClient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var logRushHttpApi = httpClient{}

type httpClient struct{}

func (c httpClient) jsonPostRequest(url string, body interface{}, responseInterface interface{}) error {
	jsonBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(jsonBody)
	response, err := http.Post(url, "application/json", requestBody)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, responseInterface)
}

func (c httpClient) RegisterStream(url, name, id, key string) (LogRushApiStreamResponse, error) {
	streamResponse := LogRushApiStreamResponse{}
	body := map[string]string{
		"alias": name,
		"id":    id,
		"key":   key,
	}

	err := c.jsonPostRequest(url+"stream/register", body, &streamResponse)
	return streamResponse, err
}

func (c httpClient) UnregisterStream(url, id, key string) (LogRushApiSuccessOrErrorResponse, error) {
	streamResponse := LogRushApiSuccessOrErrorResponse{}
	body := map[string]string{
		"id":  id,
		"key": key,
	}

	err := c.jsonPostRequest(url+"stream/unregister", body, &streamResponse)
	return streamResponse, err
}

func (c httpClient) Log(url, stream string, log LogRushLog) (LogRushApiSuccessOrErrorResponse, error) {
	streamResponse := LogRushApiSuccessOrErrorResponse{}
	body := map[string]interface{}{
		"stream":    stream,
		"log":       log.Log,
		"timestamp": log.Timestamp,
	}

	err := c.jsonPostRequest(url+"log", body, &streamResponse)
	return streamResponse, err
}

func (c httpClient) Batch(url, stream string, logs []LogRushLog) (LogRushApiSuccessOrErrorResponse, error) {
	streamResponse := LogRushApiSuccessOrErrorResponse{}
	apiLogs := []map[string]interface{}{}

	for _, log := range logs {
		apiLogs = append(apiLogs, map[string]interface{}{
			"log":       log.Log,
			"timestamp": log.Timestamp,
		})
	}

	body := map[string]interface{}{
		"stream": stream,
		"logs":   apiLogs,
	}

	err := c.jsonPostRequest(url+"batch", body, &streamResponse)
	return streamResponse, err
}
