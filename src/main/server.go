package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestHandler func(requestBody []byte) ([]byte, error)

type server struct {
	port            int
	httpServer      *http.Server
	multiplexer     *http.ServeMux
	requestHandlers map[string]map[string]RequestHandler
}

func (server *server) handleRequest(writer http.ResponseWriter, request *http.Request) {
	path, method := request.URL.Path, request.Method
	handlerFunction := server.requestHandlers[path][method]

	if handlerFunction == nil {
		http.Error(writer, fmt.Sprintf("Method '%s' is not supported.", method), http.StatusMethodNotAllowed)
		return
	}

	requestBody, error := ioutil.ReadAll(request.Body)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
		return
	}

	responseBody, error := handlerFunction(requestBody)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusBadRequest)
		return
	}

	_, error = writer.Write(responseBody)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *server) HandleRequest(path string, method string, handler RequestHandler) *server {
	log.Println(fmt.Sprintf("Registering request handler for path '%s', method '%s'...", path, method))
	if _, ok := server.requestHandlers[path]; !ok {
		server.requestHandlers[path] = map[string]RequestHandler{}
	}

	server.requestHandlers[path][method] = handler
	return server
}

func (server *server) Start() error {
	log.Println("Bootstrapping server...")
	for path, _ := range server.requestHandlers {
		server.multiplexer.HandleFunc(path, server.handleRequest)
	}

	return server.httpServer.ListenAndServe()
}

func NewServer(port int) *server {
	multiplexer := http.NewServeMux()

	return &server{
		port:            port,
		multiplexer:     multiplexer,
		requestHandlers: map[string]map[string]RequestHandler{},
		httpServer:      &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: multiplexer},
	}
}
