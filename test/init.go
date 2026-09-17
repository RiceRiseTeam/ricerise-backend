package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ricerise/internal"
	"ricerise/internal/dto"
	"time"
)

var server *httptest.Server
var client *http.Client

func init() {
	engine := internal.GetEngine()
	server = httptest.NewServer(engine)
	client = &http.Client{Timeout: 5 * time.Second}
}

func postRequest[T any](url string, request any) (int, *T) {
	return postRequestWithAuth[T](url, request, "")
}

func postRequestWithAuth[T any](url string, request any, token string) (int, *T) {
	var buf []byte = nil

	if request != nil {
		newBuf, err := json.Marshal(request)
		if err != nil {
			panic(err)
		}
		buf = newBuf
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1"+url, bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var msg dto.CommonResponse
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		panic(err)
	}

	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	if msg.Data == nil {
		return msg.Code, nil
	}

	data, err := json.Marshal(msg.Data)
	if err != nil {
		panic(err)
	}

	var result T
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}

	return msg.Code, &result
}
