package handlers

import (
	"context"
	"encoding/json"

	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/snnus/mainservice/internal/models"
)

type mainService interface {
	CreateNewSP(context.Context, models.NewServicePointRequest) (*models.ServicePoint, error)
	UpdateSP(context.Context, string, models.NewServicePointRequest) (*models.ServicePoint, error)
	DeleteSP(context.Context, string) (*models.ServicePoint, error)
	GetSPByID(context.Context, string) (*models.ServicePoint, error)
	Enqueue(context.Context, string) (*models.Ticket, error)
	Dequeue(context.Context, string) (*models.Ticket, error)
}

type MainHandler struct {
	service mainService
}

func NewMainHandler(service mainService) *MainHandler {
	return &MainHandler{service: service}
}

func (m *MainHandler) CreateNewSP(w http.ResponseWriter, r *http.Request) {
	log.Print("new service point handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newSp models.NewServicePointRequest

	if err := json.NewDecoder(r.Body).Decode(&newSp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if newSp.Name == "" || newSp.ShortName == "" || newSp.OfficeNumber == "" {
		http.Error(w, "name, shortName and officeNumber are required", http.StatusBadRequest)
		return
	}

	createdSP, err := m.service.CreateNewSP(ctx, newSp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("error creating service point: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdSP); err != nil {
		log.Printf("failed to encode response: %s", err)
	}

	log.Printf("201 created - service point ID: %d", createdSP.ID)
}

func (m *MainHandler) UpdateSP(w http.ResponseWriter, r *http.Request) {
	log.Print("update service point handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newSp models.NewServicePointRequest

	vars := mux.Vars(r)
	id := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&newSp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if newSp.Name == "" || newSp.ShortName == "" || newSp.OfficeNumber == "" {
		http.Error(w, "name, shortName and officeNumber are required", http.StatusBadRequest)
		return
	}

	updatedSP, err := m.service.UpdateSP(ctx, id, newSp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("error updating service point: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(updatedSP); err != nil {
		log.Printf("failed to encode response: %s", err)
	}

	log.Printf("200 ok - service point ID: %d", updatedSP.ID)
}

func (m *MainHandler) DeleteSP(w http.ResponseWriter, r *http.Request) {
	log.Print("delete service point handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	deletedSP, err := m.service.DeleteSP(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("error deleting service point: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(deletedSP); err != nil {
		log.Printf("failed to encode response: %s", err)
	}

	log.Printf("200 ok - service point ID: %d", deletedSP.ID)
}

func (m *MainHandler) GetSP(w http.ResponseWriter, r *http.Request) {
	log.Print("get service point handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	sp, err := m.service.GetSPByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("error getting service point: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(sp); err != nil {
		log.Printf("failed to encode response: %s", err)
	}

	log.Printf("200 ok - service point ID: %d", sp.ID)
}

func (m *MainHandler) Enqueue(w http.ResponseWriter, r *http.Request) {
	log.Print("enqueue handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	ticket, err := m.service.Enqueue(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("error getting service point: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		log.Printf("failed to encode response: %s", err)
	}

	log.Printf("200 ok")
}

func (m *MainHandler) Dequeue(w http.ResponseWriter, r *http.Request) {
	log.Print("enqueue handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	ticket, err := m.service.Dequeue(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("error getting service point: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		log.Printf("failed to encode response: %s", err)
	}

	log.Printf("200 ok")
}
