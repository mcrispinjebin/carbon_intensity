package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"carbon_intensity/adapter"
	"carbon_intensity/handler"
	"carbon_intensity/processor"
)

func main() {
	r := mux.NewRouter()

	client := adapter.NewHTTPClient()
	processorService := processor.Processor{
		DataClient: client,
	}

	apiHandler := handler.Handler{
		Processor: &processorService,
	}

	r.HandleFunc("/slots", apiHandler.GetSlotsHandler).Methods("GET")

	log.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
