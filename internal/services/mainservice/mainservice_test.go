package mainservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/snnus/mainservice/internal/models"
	"github.com/snnus/mainservice/internal/services/mainservice"
	"github.com/snnus/mainservice/internal/services/mainservice/mocks"
	"github.com/stretchr/testify/suite"
)

type MainServiceTestSuite struct {
	suite.Suite
	ctx        context.Context
	storage    *mocks.MockMainStorage
	httpClient *mocks.MockMainClient
	producer   *mocks.MockMainProducer
	service    *mainservice.MainService
}

func (s *MainServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = mocks.NewMockMainStorage(s.T())
	s.httpClient = mocks.NewMockMainClient(s.T())
	s.producer = mocks.NewMockMainProducer(s.T())
	s.service = mainservice.NewMainService(s.storage, s.httpClient, s.producer)
}

func TestMainServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MainServiceTestSuite))
}

func (s *MainServiceTestSuite) TestCreateNewSP_Success() {
	spRequest := models.NewServicePointRequest{
		Name:         "Test Service Point",
		ShortName:    "TSP",
		OfficeNumber: "101",
	}
	expectedSP := &models.ServicePoint{
		ID:           123,
		Name:         "Test Service Point",
		ShortName:    "TSP",
		OfficeNumber: "101",
	}

	s.storage.On("CreateServicePoint", s.ctx, spRequest).Return(expectedSP, nil)

	result, err := s.service.CreateNewSP(s.ctx, spRequest)

	s.NoError(err)
	s.Equal(expectedSP, result)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestCreateNewSP_MissingName() {
	spRequest := models.NewServicePointRequest{
		ShortName:    "TSP",
		OfficeNumber: "101",
	}

	result, err := s.service.CreateNewSP(s.ctx, spRequest)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "name is required")
	s.storage.AssertNotCalled(s.T(), "CreateServicePoint")
}

func (s *MainServiceTestSuite) TestCreateNewSP_MissingShortName() {
	spRequest := models.NewServicePointRequest{
		Name:         "Test Service Point",
		OfficeNumber: "101",
	}

	result, err := s.service.CreateNewSP(s.ctx, spRequest)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "short name is required")
	s.storage.AssertNotCalled(s.T(), "CreateServicePoint")
}

func (s *MainServiceTestSuite) TestCreateNewSP_MissingOfficeNumber() {
	spRequest := models.NewServicePointRequest{
		Name:      "Test Service Point",
		ShortName: "TSP",
	}

	result, err := s.service.CreateNewSP(s.ctx, spRequest)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "office number is required")
	s.storage.AssertNotCalled(s.T(), "CreateServicePoint")
}

func (s *MainServiceTestSuite) TestCreateNewSP_StorageError() {
	spRequest := models.NewServicePointRequest{
		Name:         "Test Service Point",
		ShortName:    "TSP",
		OfficeNumber: "101",
	}
	storageError := errors.New("database error")

	s.storage.On("CreateServicePoint", s.ctx, spRequest).Return(nil, storageError)

	result, err := s.service.CreateNewSP(s.ctx, spRequest)

	s.Error(err)
	s.Nil(result)
	s.Equal(storageError, err)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestUpdateSP_Success() {
	id := "123"
	spRequest := models.NewServicePointRequest{
		Name:         "Updated Service Point",
		ShortName:    "USP",
		OfficeNumber: "202",
	}
	expectedSP := &models.ServicePoint{
		ID:           123,
		Name:         "Updated Service Point",
		ShortName:    "USP",
		OfficeNumber: "202",
	}

	s.storage.On("UpdateServicePoint", s.ctx, id, spRequest).Return(expectedSP, nil)

	result, err := s.service.UpdateSP(s.ctx, id, spRequest)

	s.NoError(err)
	s.Equal(expectedSP, result)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestUpdateSP_ValidationErrors() {
	testCases := []struct {
		name     string
		request  models.NewServicePointRequest
		expected string
	}{
		{
			name: "missing name",
			request: models.NewServicePointRequest{
				ShortName:    "USP",
				OfficeNumber: "202",
			},
			expected: "name is required",
		},
		{
			name: "missing short name",
			request: models.NewServicePointRequest{
				Name:         "Updated Service Point",
				OfficeNumber: "202",
			},
			expected: "short name is required",
		},
		{
			name: "missing office number",
			request: models.NewServicePointRequest{
				Name:      "Updated Service Point",
				ShortName: "USP",
			},
			expected: "office number is required",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result, err := s.service.UpdateSP(s.ctx, "123", tc.request)
			s.Error(err)
			s.Nil(result)
			s.Contains(err.Error(), tc.expected)
			s.storage.AssertNotCalled(s.T(), "UpdateServicePoint")
		})
	}
}

func (s *MainServiceTestSuite) TestDeleteSP_Success() {
	id := "123"
	expectedSP := &models.ServicePoint{
		ID:           123,
		Name:         "Service Point",
		ShortName:    "SP",
		OfficeNumber: "101",
	}

	s.storage.On("DeleteServicePoint", s.ctx, id).Return(expectedSP, nil)

	result, err := s.service.DeleteSP(s.ctx, id)

	s.NoError(err)
	s.Equal(expectedSP, result)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestDeleteSP_StorageError() {
	id := "123"
	storageError := errors.New("delete error")

	s.storage.On("DeleteServicePoint", s.ctx, id).Return(nil, storageError)

	result, err := s.service.DeleteSP(s.ctx, id)

	s.Error(err)
	s.Nil(result)
	s.Equal(storageError, err)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestGetSPByID_Success() {
	id := "123"
	expectedSP := &models.ServicePoint{
		ID:           123,
		Name:         "Service Point",
		ShortName:    "SP",
		OfficeNumber: "101",
	}

	s.storage.On("GetServicePointByID", s.ctx, id).Return(expectedSP, nil)

	result, err := s.service.GetSPByID(s.ctx, id)

	s.NoError(err)
	s.Equal(expectedSP, result)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestGetSPByID_StorageError() {
	id := "123"
	storageError := errors.New("not found")

	s.storage.On("GetServicePointByID", s.ctx, id).Return(nil, storageError)

	result, err := s.service.GetSPByID(s.ctx, id)

	s.Error(err)
	s.Nil(result)
	s.Equal(storageError, err)
	s.storage.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestEnqueue_Success() {
	id := "123"
	shortName := "SP"
	expectedTicket := &models.Ticket{
		Ticket: "T001",
	}

	s.storage.On("GetShortNameById", s.ctx, id).Return(shortName, nil)
	s.httpClient.On("Enqueue", s.ctx, id, shortName).Return(expectedTicket, nil)

	result, err := s.service.Enqueue(s.ctx, id)

	s.NoError(err)
	s.Equal(expectedTicket, result)
	s.storage.AssertExpectations(s.T())
	s.httpClient.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestEnqueue_GetShortNameError() {
	id := "123"
	storageError := errors.New("service point not found")

	s.storage.On("GetShortNameById", s.ctx, id).Return("", storageError)

	result, err := s.service.Enqueue(s.ctx, id)

	s.Error(err)
	s.Nil(result)
	s.Equal(storageError, err)
	s.storage.AssertExpectations(s.T())
	s.httpClient.AssertNotCalled(s.T(), "Enqueue")
}

func (s *MainServiceTestSuite) TestEnqueue_HttpClientError() {
	id := "123"
	shortName := "SP"
	clientError := errors.New("enqueue failed")

	s.storage.On("GetShortNameById", s.ctx, id).Return(shortName, nil)
	s.httpClient.On("Enqueue", s.ctx, id, shortName).Return(nil, clientError)

	result, err := s.service.Enqueue(s.ctx, id)

	s.Error(err)
	s.Nil(result)
	s.Equal(clientError, err)
	s.storage.AssertExpectations(s.T())
	s.httpClient.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestDequeue_Success() {
	id := "123"
	officeNumber := "101"
	expectedTicket := &models.Ticket{
		Ticket: "T001",
	}

	s.httpClient.On("Dequeue", s.ctx, id).Return(expectedTicket, nil)
	s.storage.On("GetOfficeNumberById", s.ctx, id).Return(officeNumber, nil)
	s.producer.On("PublishTicket", s.ctx, expectedTicket.Ticket, officeNumber).Return(nil)

	result, err := s.service.Dequeue(s.ctx, id)

	s.NoError(err)
	s.Equal(expectedTicket, result)
	s.httpClient.AssertExpectations(s.T())
	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestDequeue_DequeueError() {
	id := "123"
	clientError := errors.New("dequeue failed")

	s.httpClient.On("Dequeue", s.ctx, id).Return(nil, clientError)

	result, err := s.service.Dequeue(s.ctx, id)

	s.Error(err)
	s.Nil(result)
	s.Equal(clientError, err)
	s.httpClient.AssertExpectations(s.T())
	s.storage.AssertNotCalled(s.T(), "GetOfficeNumberById")
	s.producer.AssertNotCalled(s.T(), "PublishTicket")
}

func (s *MainServiceTestSuite) TestDequeue_GetOfficeNumberError() {
	id := "123"
	expectedTicket := &models.Ticket{
		Ticket: "T001",
	}
	storageError := errors.New("office number not found")

	s.httpClient.On("Dequeue", s.ctx, id).Return(expectedTicket, nil)
	s.storage.On("GetOfficeNumberById", s.ctx, id).Return("", storageError)

	result, err := s.service.Dequeue(s.ctx, id)

	s.Error(err)
	s.Nil(result)
	s.Equal(storageError, err)
	s.httpClient.AssertExpectations(s.T())
	s.storage.AssertExpectations(s.T())
	s.producer.AssertNotCalled(s.T(), "PublishTicket")
}

func (s *MainServiceTestSuite) TestDequeue_PublishError() {
	id := "123"
	officeNumber := "101"
	expectedTicket := &models.Ticket{
		Ticket: "T001",
	}
	publishError := errors.New("publish failed")

	s.httpClient.On("Dequeue", s.ctx, id).Return(expectedTicket, nil)
	s.storage.On("GetOfficeNumberById", s.ctx, id).Return(officeNumber, nil)
	s.producer.On("PublishTicket", s.ctx, expectedTicket.Ticket, officeNumber).Return(publishError)

	// Note: Publish error is logged but doesn't cause the method to fail
	result, err := s.service.Dequeue(s.ctx, id)

	s.NoError(err) // Publish error is only logged, not returned
	s.Equal(expectedTicket, result)
	s.httpClient.AssertExpectations(s.T())
	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestDequeue_PublishErrorWithLogging() {
	id := "123"
	officeNumber := "101"
	expectedTicket := &models.Ticket{
		Ticket: "T001",
	}
	publishError := errors.New("publish failed")

	// Mock the log.Printf to verify it's called
	// Since we can't easily intercept log.Printf without modifying the service,
	// we'll just verify the method still succeeds despite the error
	s.httpClient.On("Dequeue", s.ctx, id).Return(expectedTicket, nil)
	s.storage.On("GetOfficeNumberById", s.ctx, id).Return(officeNumber, nil)
	s.producer.On("PublishTicket", s.ctx, expectedTicket.Ticket, officeNumber).Return(publishError)

	result, err := s.service.Dequeue(s.ctx, id)

	s.NoError(err)
	s.Equal(expectedTicket, result)
	// The error is logged internally but doesn't propagate
	s.producer.AssertExpectations(s.T())
}

func (s *MainServiceTestSuite) TestIntegration_EnqueueThenDequeue() {
	id := "123"
	shortName := "SP"
	officeNumber := "101"
	ticket := &models.Ticket{
		Ticket: "T001",
	}

	// Enqueue setup
	s.storage.On("GetShortNameById", s.ctx, id).Return(shortName, nil)
	s.httpClient.On("Enqueue", s.ctx, id, shortName).Return(ticket, nil)

	// Dequeue setup - note we need to clear previous expectations
	// or use mock.Once for the specific test
	s.SetupTest()

	s.httpClient.On("Dequeue", s.ctx, id).Return(ticket, nil)
	s.storage.On("GetOfficeNumberById", s.ctx, id).Return(officeNumber, nil)
	s.producer.On("PublishTicket", s.ctx, ticket.Ticket, officeNumber).Return(nil)

	enqueueResult, enqueueErr := s.service.Enqueue(s.ctx, id)
	s.NoError(enqueueErr)
	s.Equal(ticket, enqueueResult)

	dequeueResult, dequeueErr := s.service.Dequeue(s.ctx, id)
	s.NoError(dequeueErr)
	s.Equal(ticket, dequeueResult)
}

func (s *MainServiceTestSuite) TestNilChecks() {
	// Test that service methods handle edge cases appropriately
	s.Run("CreateNewSP with empty context", func() {
		spRequest := models.NewServicePointRequest{
			Name:         "Test",
			ShortName:    "T",
			OfficeNumber: "101",
		}
		result, err := s.service.CreateNewSP(nil, spRequest)
		s.Error(err)
		s.Nil(result)
	})

	s.Run("GetSPByID with empty string", func() {
		result, err := s.service.GetSPByID(s.ctx, "")
		s.Error(err)
		s.Nil(result)
	})
}
