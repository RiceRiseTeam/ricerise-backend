package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ricerise/internal"
	"time"
)

var server *httptest.Server
var client *http.Client

func init() {
	engine := internal.GetEngine()
	server = httptest.NewServer(engine)
	client = &http.Client{Timeout: 5 * time.Second}
}

func postRequest(url string, dto any) (*http.Response, map[string]any) {
	buf, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	resp, err := client.Post(server.URL+"/api/v1"+url, "application/json", bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	var msg map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		panic(err)
	}

	return resp, msg
}
