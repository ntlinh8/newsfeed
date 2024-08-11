package main

import (
	"log"

	"newsfeed/pkg/handler/http"
)

func main() {

	// TODO: read from env
	httpConfig := http.HttpHandlerConfig{UserGrpcAddr: "localhost:44444"}

	httpHandler, err := http.New(httpConfig)
	if err != nil {
		log.Println("failed to create http handler", err)
		return
	}

	httpHandler.Start() // block

	// httpHandler.Stop()
}
