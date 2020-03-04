package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	PathWord                   = "/word"
	PathHistory                = "/history"
	PathSentence               = "/sentence"
	PropertyKeyHistory         = "history"
	PropertyKeyGopherWord      = "gopher-word"
	PropertyKeyEnglishWord     = "english-word"
	PropertyKeyGopherSentence  = "gopher-sentence"
	PropertyKeyEnglishSentence = "english-sentence"
)

var worklog = NewWorklogManager(100)
var wordTranslatorInstance = NewWordTranslator()
var sentenceTranslatorInstance = NewSentenceTranslator()
var port = flag.Int("port", 80, "Port for the HTTP server to listen to")

func listWorklog([]byte) ([]byte, error) {
	log.Println("Listing worklog history...")
	return json.Marshal(map[string][]WorklogEntry{PropertyKeyHistory: worklog.GetAll()})
}

func translate(rawBody []byte, inputKey string, outputKey string, translator Translator) ([]byte, error) {
	var request map[string]string
	error := json.Unmarshal(rawBody, &request)
	if error != nil {
		return nil, error
	}

	output, error := translator.Translate(request[inputKey])
	if error != nil {
		return nil, error
	}

	log.Println(fmt.Sprintf("Translated input '%s' into '%s'.", request[inputKey], output))
	worklog.Persist(map[string]string{request[inputKey]: output})
	return json.Marshal(map[string]string{outputKey: output})
}

func translateWord(rawBody []byte) ([]byte, error) {
	return translate(rawBody, PropertyKeyEnglishWord, PropertyKeyGopherWord, wordTranslatorInstance)
}

func translateSentence(rawBody []byte) ([]byte, error) {
	return translate(rawBody, PropertyKeyEnglishSentence, PropertyKeyGopherSentence, sentenceTranslatorInstance)
}

func main() {
	flag.Parse()

	log.Println(fmt.Sprintf("Starting server on port %d...", *port))

	err := NewServer(*port).
		HandleRequest(PathHistory, http.MethodGet, listWorklog).
		HandleRequest(PathWord, http.MethodPost, translateWord).
		HandleRequest(PathSentence, http.MethodPost, translateSentence).
		Start()

	if err != nil {
		log.Fatal(err)
	}
}
