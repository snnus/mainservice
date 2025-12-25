package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/snnus/mainservice/config"
	"github.com/snnus/mainservice/internal/client"
	"github.com/snnus/mainservice/internal/handlers"
	"github.com/snnus/mainservice/internal/producer"
	"github.com/snnus/mainservice/internal/services/mainservice"
	"github.com/snnus/mainservice/internal/storage/spstorage"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("CONFIG_PATH"))

	if err != nil {
		panic(err)
	}

	spClient := client.NewClient(cfg)
	spProducer := producer.NewSPProducer(cfg)

	spStorage, close, err := spstorage.NewSPStorage(cfg)

	if err != nil {
		panic(err)
	}
	defer close()

	spService := mainservice.NewSPService(spStorage, spClient, spProducer)
	spHandler := handlers.NewSPHandler(spService)

	r := mux.NewRouter()

	r.HandleFunc("/servicepoint/{id:[0-9]+}", spHandler.UpsertSP).Methods("PUT", "POST")
	r.HandleFunc("/servicepoint/{id:[0-9]+}", spHandler.GetSP).Methods("GET")
	r.HandleFunc("/servicepoint/{id:[0-9]+}", spHandler.DeleteSP).Methods("DELETE")
	r.HandleFunc("/enqueue/{id:[0-9]+}", spHandler.Enqueue).Methods("POST")
	r.HandleFunc("/dequeue/{id:[0-9]+}", spHandler.Dequeue).Methods("POST")

	log.Print("listening now")

	http.ListenAndServe("0.0.0.0:8080", r)
}
