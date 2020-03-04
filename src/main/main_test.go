package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestThatWhenTranslateWordIsCalled_ThenCorrectResponseIsReturned(t *testing.T) {
	testResponse := map[string]string{}
	testRequestPayload, _ := json.Marshal(map[string]string{"english-word": "apple"})

	rawResponse, error := translateWord(testRequestPayload)
	if error != nil {
		t.Errorf("Unexpected error %s.", error.Error())
	}

	json.Unmarshal(rawResponse, &testResponse)

	if expectedResponse := map[string]string{"gopher-word": "gapple"}; !reflect.DeepEqual(testResponse, expectedResponse) {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", testResponse, expectedResponse)
	}
}

func TestThatWordIsNotValidWhenTranslateWordIsCalled_ThenCorrectErrorIsReturned(t *testing.T) {
	testRequestPayloadEmptyWord, _ := json.Marshal(map[string]string{"english-word": ""})
	testRequestPayloadInvalidSymbol, _ := json.Marshal(map[string]string{"english-word": "'"})

	_, error := translateWord(testRequestPayloadEmptyWord)
	if error == nil {
		t.Errorf("Expected error for an empty word, got nil.")
	}

	_, error = translateWord(testRequestPayloadInvalidSymbol)
	if error == nil {
		t.Errorf("Expected error for an invalid symbol, got nil.")
	}
}

func TestThatWhenTranslateSentenceIsCalled_ThenCorrectResponseIsReturned(t *testing.T) {
	testResponse := map[string]string{}
	testRequestPayload, _ := json.Marshal(map[string]string{"english-sentence": "square, apple, xray, chair!"})

	rawResponse, error := translateSentence(testRequestPayload)
	if error != nil {
		t.Errorf("Unexpected error %s.", error.Error())
	}

	json.Unmarshal(rawResponse, &testResponse)

	if expectedResponse := map[string]string{"gopher-sentence": "aresquogo, gapple, gexray, airchogo!"}; !reflect.DeepEqual(testResponse, expectedResponse) {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", testResponse, expectedResponse)
	}
}

func TestThatSentenceIsNotValidWhenTranslateSentenceIsCalled_ThenCorrectErrorIsReturned(t *testing.T) {
	testRequestPayloadEmptyWord, _ := json.Marshal(map[string]string{"english-sentence": ""})
	testRequestPayloadInvalidSymbol, _ := json.Marshal(map[string]string{"english-sentence": "'"})

	_, error := translateSentence(testRequestPayloadEmptyWord)
	if error == nil {
		t.Errorf("Expected error for an empty word, got nil.")
	}

	_, error = translateSentence(testRequestPayloadInvalidSymbol)
	if error == nil {
		t.Errorf("Expected error for an invalid symbol, got nil.")
	}
}
