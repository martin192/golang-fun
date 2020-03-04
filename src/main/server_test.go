package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhenRequestHandlersAreConfigured_ThenHandlersAreStoredAccordingly(t *testing.T) {
	testPath := "/echo"
	testHandler := func(body []byte) ([]byte, error) { return nil, nil }
	testServer := NewServer(1234).
		HandleRequest(testPath, http.MethodGet, testHandler).
		HandleRequest(testPath, http.MethodPut, testHandler).
		HandleRequest(testPath, http.MethodPost, testHandler)

	if _, exists := testServer.requestHandlers[testPath][http.MethodGet]; !exists {
		t.Errorf("No handler found for path %s and method %s.", testPath, http.MethodGet)
	}

	if _, exists := testServer.requestHandlers[testPath][http.MethodPut]; !exists {
		t.Errorf("No handler found for path %s and method %s.", testPath, http.MethodPut)
	}

	if _, exists := testServer.requestHandlers[testPath][http.MethodPost]; !exists {
		t.Errorf("No handler found for path %s and method %s.", testPath, http.MethodPost)
	}
}

func TestWhenCorrectRequestIsSent_ThenCorrectResponseIsReturned(t *testing.T) {
	testPayload := "stuff"
	requestRecorder := httptest.NewRecorder()
	testServerPort := 1024 + (rand.Int() % 10000)
	testRequest := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBuffer([]byte(testPayload)))
	testServer := NewServer(testServerPort).
		HandleRequest("/echo", http.MethodPost, func(body []byte) ([]byte, error) { return body, nil })

	testServer.handleRequest(requestRecorder, testRequest)

	if foundCode := requestRecorder.Result().StatusCode; foundCode != http.StatusOK {
		t.Errorf("Unexpected status code, got: %d, wanted: %d.", foundCode, http.StatusOK)
	}

	if responseBody, _ := ioutil.ReadAll(requestRecorder.Result().Body); string(responseBody) != testPayload {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", string(responseBody), testPayload)
	}
}

func TestWhenRequestMethodIsNotValid_ThenCorrectResponseIsReturned(t *testing.T) {
	requestRecorder := httptest.NewRecorder()
	testServerPort := 1024 + (rand.Int() % 10000)
	testRequest := httptest.NewRequest(http.MethodGet, "/echo", nil)
	testServer := NewServer(testServerPort).
		HandleRequest("/echo", http.MethodPost, func(body []byte) ([]byte, error) { return body, nil })

	testServer.handleRequest(requestRecorder, testRequest)

	if foundCode := requestRecorder.Result().StatusCode; foundCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected status code, got: %d, wanted: %d.", foundCode, http.StatusMethodNotAllowed)
	}
}

func TestWhenRequestProcessingFails_ThenCorrectResponseIsReturned(t *testing.T) {
	requestRecorder := httptest.NewRecorder()
	testServerPort := 1024 + (rand.Int() % 10000)
	testRequest := httptest.NewRequest(http.MethodPost, "/echo", nil)
	testServer := NewServer(testServerPort).
		HandleRequest("/echo", http.MethodPost, func(body []byte) ([]byte, error) { return nil, errors.New("error") })

	testServer.handleRequest(requestRecorder, testRequest)

	if foundCode := requestRecorder.Result().StatusCode; foundCode != http.StatusBadRequest {
		t.Errorf("Unexpected status code, got: %d, wanted: %d.", foundCode, http.StatusBadRequest)
	}
}
