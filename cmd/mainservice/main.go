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
	"github.com/snnus/mainservice/internal/storage/pgstorage"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("CONFIG_PATH"))

	if err != nil {
		panic(err)
	}

	db, err := pgstorage.NewConnection(*cfg)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	mainClient := client.NewClient(cfg)
	mainProducer := producer.NewKafkaProducer(cfg)

	mainStorage := pgstorage.NewPGStorage(db)
	mainService := mainservice.NewMainService(mainStorage, mainClient, mainProducer)
	mainHandler := handlers.NewMainHandler(mainService)

	r := mux.NewRouter()

	r.HandleFunc("/servicepoint", mainHandler.CreateNewSP).Methods("POST")
	r.HandleFunc("/servicepoint/{id:[0-9]+}", mainHandler.UpdateSP).Methods("PUT")
	r.HandleFunc("/servicepoint/{id:[0-9]+}", mainHandler.GetSP).Methods("GET")
	r.HandleFunc("/servicepoint/{id:[0-9]+}", mainHandler.DeleteSP).Methods("DELETE")
	r.HandleFunc("/enqueue/{id:[0-9]+}", mainHandler.Enqueue).Methods("POST")
	r.HandleFunc("/dequeue/{id:[0-9]+}", mainHandler.Dequeue).Methods("POST")

	log.Print("listening now")

	http.ListenAndServe("0.0.0.0:8080", r)
}
