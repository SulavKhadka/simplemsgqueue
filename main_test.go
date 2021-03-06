package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sulavkhadka/queue"
)

func TestCreateTopic(t *testing.T) {

	srv := newServer()

	samplePayload := []byte(`{"topic_name": "Test","queue": ""}`)

	//API endpoint setup
	req, err := http.NewRequest("PUT", "/topic", bytes.NewBuffer(samplePayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.CreateTopic)

	//Making the API call
	handler.ServeHTTP(rr, req)

	if req.Method != "PUT" {
		t.Errorf("Expected 'PUT' request, got '%s'", req.Method)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedResponse := struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Err    string `json:"err"`
	}{"Success", "Topic Test Added.", ""}

	actualResponse := struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Err    string `json:"err"`
	}{}
	actualPayload := json.NewDecoder(rr.Body)
	err = actualPayload.Decode(&actualResponse)
	check(err)

	if actualResponse != expectedResponse {
		t.Errorf("newTopic was incorrect, got: %s (%s), want: %s (%s).",
			expectedResponse, reflect.TypeOf(rr.Body.String()), rr.Body.String(), reflect.TypeOf(expectedResponse))
	}

}

func TestSendMessage(t *testing.T) {

	srv := newServer()

	// FIXME: Rethink the way to make the create topic call.
	newTopic := []byte(`{"topic_name": "TestTopic","queue": ""}`)

	//API endpoint setup
	req, err := http.NewRequest("PUT", "/topic", bytes.NewBuffer(newTopic))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.CreateTopic)

	//Making the API call
	handler.ServeHTTP(rr, req)

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ //

	//Sample payload setup by converting struct to a byte array
	samplePayload, err := json.Marshal(message{"message my dude", 0, "TestTopic"})
	check(err)

	//API endpoint setup
	req, err = http.NewRequest("POST", "/topic", bytes.NewBuffer(samplePayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(srv.SendMessage)

	//Making the API call
	handler.ServeHTTP(rr, req)

	if req.Method != "POST" {
		t.Errorf("Expected 'POST' request, got '%s'", req.Method)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedResponse := struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Err    string `json:"err"`
	}{"Success", "Recieved Msg message my dude // ID: 0 // Topic: TestTopic", ""}

	actualResponse := struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Err    string `json:"err"`
	}{}
	actualPayload := json.NewDecoder(rr.Body)
	err = actualPayload.Decode(&actualResponse)
	check(err)

	if actualResponse != expectedResponse {
		t.Errorf("send message was incorrect, got: %s (%s), want: %s (%s).",
			actualResponse, reflect.TypeOf(actualResponse), expectedResponse, reflect.TypeOf(expectedResponse))
	}

}

func TestGetMessage(t *testing.T) {

	srv := newServer()

	// FIXME: Rethink the way to make the create topic call.
	newTopic := []byte(`{"topic_name": "TestTopic","queue": ""}`)

	//API endpoint setup
	req, err := http.NewRequest("PUT", "/topic", bytes.NewBuffer(newTopic))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.CreateTopic)

	//Making the API call
	handler.ServeHTTP(rr, req)

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ //

	//Sample payload setup by converting struct to a byte array
	sampleMessage, err := json.Marshal(message{"message my dude", 0, "TestTopic"})
	check(err)

	//API endpoint setup
	req, err = http.NewRequest("POST", "/topic", bytes.NewBuffer(sampleMessage))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(srv.SendMessage)

	//Making the API call
	handler.ServeHTTP(rr, req)

	// FIXME: fix this mess above of creating topics and populating them

	//Sample payload setup by converting struct to a byte array
	samplePayload, err := json.Marshal(struct {
		TopicName string `json:"topic_name"`
	}{"TestTopic"})
	check(err)

	//API endpoint setup
	req, err = http.NewRequest("GET", "/topic", bytes.NewBuffer(samplePayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(srv.GetMessage)

	//Making the API call
	handler.ServeHTTP(rr, req)

	if req.Method != "GET" {
		t.Errorf("Expected 'GET' request, got '%s'", req.Method)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedResponse := struct {
		Topic string     `json:"topic"`
		Msg   queue.Item `json:"msg"`
		Err   string     `json:"err"`
	}{"TestTopic", queue.Item{1, "message my dude", "", "TestTopic"}, ""}

	actualResponse := struct {
		Topic string     `json:"topic"`
		Msg   queue.Item `json:"msg"`
		Err   string     `json:"err"`
	}{}
	actualPayload := json.NewDecoder(rr.Body)
	err = actualPayload.Decode(&actualResponse)
	check(err)

	actualResponse.Msg.TimeStamp = "" //FIXME: This is to circumvent the timestamp check. need a better way to do so.

	if actualResponse != expectedResponse {
		t.Errorf("send message was incorrect, got: %v (%s), want: %v (%s).",
			actualResponse, reflect.TypeOf(actualResponse), expectedResponse, reflect.TypeOf(expectedResponse))
	}

}
