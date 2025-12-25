package spservice

import (
	"context"
	"fmt"
	"log"

	"github.com/snnus/mainservice/internal/models"
)

type SPStorage interface {
	UpsertServicePoint(ctx context.Context, id string, sp models.NewServicePointRequest) (*models.ServicePoint, error)
	DeleteServicePoint(ctx context.Context, id string) (*models.ServicePoint, error)
	GetServicePointByID(ctx context.Context, id string) (*models.ServicePoint, error)
	GetShortNameById(ctx context.Context, is string) (string, error)
	GetOfficeNumberById(ctx context.Context, is string) (string, error)
}

type SPClient interface {
	Enqueue(ctx context.Context, id string, shortname string) (*models.Ticket, error)
	Dequeue(ctx context.Context, id string) (*models.Ticket, error)
}

type SPProducer interface {
	PublishTicket(ctx context.Context, ticket, officeNumber string) error
}

type SPService struct {
	storage    SPStorage
	httpClient SPClient
	producer   SPProducer
}

func NewSPService(storage SPStorage, httpClient SPClient, producer SPProducer) *SPService {
	return &SPService{storage: storage, httpClient: httpClient, producer: producer}
}

func (m *SPService) UpsertSP(ctx context.Context, id string, sp models.NewServicePointRequest) (*models.ServicePoint, error) {
	if sp.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if sp.ShortName == "" {
		return nil, fmt.Errorf("short name is required")
	}
	if sp.OfficeNumber == "" {
		return nil, fmt.Errorf("office number is required")
	}
	updatedSP, err := m.storage.UpsertServicePoint(ctx, id, sp)
	if err != nil {
		return nil, err
	}
	return updatedSP, err
}

func (m *SPService) DeleteSP(ctx context.Context, id string) (*models.ServicePoint, error) {
	deletedSP, err := m.storage.DeleteServicePoint(ctx, id)
	if err != nil {
		return nil, err
	}
	return deletedSP, err
}

func (m *SPService) GetSPByID(ctx context.Context, id string) (*models.ServicePoint, error) {
	sp, err := m.storage.GetServicePointByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return sp, err
}

// func (m *MainService) GetAllSP(ctx context.Context) (*[]models.ServicePoint, error) {

// }

func (m *SPService) Enqueue(ctx context.Context, id string) (*models.Ticket, error) {
	shortName, err := m.storage.GetShortNameById(ctx, id)
	if err != nil {
		return nil, err
	}

	ticket, err := m.httpClient.Enqueue(ctx, id, shortName)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (m *SPService) Dequeue(ctx context.Context, id string) (*models.Ticket, error) {
	ticket, err := m.httpClient.Dequeue(ctx, id)
	if err != nil {
		return nil, err
	}

	officeNumber, err := m.storage.GetOfficeNumberById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = m.producer.PublishTicket(ctx, ticket.Ticket, officeNumber)

	if err != nil {
		log.Printf("%s", err.Error())
	}

	return ticket, nil
}
