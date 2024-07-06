package worker

import (
	"context"
	"github.com/matnich89/national-rail-client/nationalrail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
	"trainstats-scraper/model"
)

type MockRailClient struct {
	mock.Mock
}

func (m *MockRailClient) GetDepartures(code nationalrail.CRSCode, opts ...nationalrail.RequestOption) (*nationalrail.StationBoard, error) {
	args := m.Called(code, opts)
	return args.Get(0).(*nationalrail.StationBoard), args.Error(1)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) PushToQueue(ctx context.Context, id model.DepartingTrainId) {
	m.Called(ctx, id)
}

func TestWorker_Work(t *testing.T) {
	mockRailClient := new(MockRailClient)
	mockRedisClient := new(MockRedisClient)

	serviceChan := make(chan model.DepartingTrainId, 10)
	worker := &Worker{
		ID:           1,
		Stations:     []nationalrail.CRSCode{"STA"},
		ServiceChan:  serviceChan,
		NRClient:     mockRailClient,
		InitialDelay: 0,
		Ticker:       time.NewTicker(100 * time.Millisecond),
		RedisClient:  mockRedisClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	scheduledTime := "12:00"
	mockRailClient.On("GetDepartures", nationalrail.CRSCode("STA"), mock.Anything).Return(&nationalrail.StationBoard{
		Services: []*nationalrail.Service{
			{
				ID:                       "SERVICE1",
				ScheduledTimeOfDeparture: &scheduledTime,
			},
		},
	}, nil)

	go worker.Work(ctx, wg)

	// Wait for the worker to finish
	wg.Wait()

	// Check that the service was sent to the channel
	select {
	case train := <-serviceChan:
		assert.Equal(t, "SERVICE1", train.ID)
	default:
		t.Error("Expected a train in the channel, but found none")
	}

	mockRailClient.AssertExpectations(t)
}

func TestWorker_checkStation(t *testing.T) {
	mockRailClient := new(MockRailClient)
	serviceChan := make(chan model.DepartingTrainId, 10)

	worker := &Worker{
		ID:          1,
		ServiceChan: serviceChan,
		NRClient:    mockRailClient,
	}

	scheduledTime1 := "12:00"
	scheduledTime2 := "13:00"
	mockRailClient.On("GetDepartures", nationalrail.CRSCode("STA"), mock.Anything).Return(&nationalrail.StationBoard{
		Services: []*nationalrail.Service{
			{
				ID:                       "SERVICE1",
				ScheduledTimeOfDeparture: &scheduledTime1,
			},
			{
				ID:                       "SERVICE2",
				ScheduledTimeOfDeparture: &scheduledTime2,
			},
		},
	}, nil)

	err := worker.checkStation(context.Background(), "STA")

	assert.NoError(t, err)
	assert.Len(t, serviceChan, 2)

	train1 := <-serviceChan
	assert.Equal(t, "SERVICE1", train1.ID)

	train2 := <-serviceChan
	assert.Equal(t, "SERVICE2", train2.ID)

	mockRailClient.AssertExpectations(t)
}

func TestWorker_processService(t *testing.T) {
	serviceChan := make(chan model.DepartingTrainId, 1)
	worker := &Worker{
		ID:          1,
		ServiceChan: serviceChan,
	}

	scheduledTime := "12:00"
	service := &nationalrail.Service{
		ID:                       "SERVICE1",
		ScheduledTimeOfDeparture: &scheduledTime,
	}

	err := worker.processService(context.Background(), service)

	assert.NoError(t, err)
	assert.Len(t, serviceChan, 1)

	train := <-serviceChan
	assert.Equal(t, "SERVICE1", train.ID)
}
