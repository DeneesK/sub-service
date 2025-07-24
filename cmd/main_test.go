package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/DeneesK/sub-service/internal/model"
	"github.com/DeneesK/sub-service/internal/router"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockSubscriptionService struct {
	data map[string]*model.Subscription
	mu   sync.RWMutex
}

func NewMockSubscriptionService() *MockSubscriptionService {
	return &MockSubscriptionService{
		data: make(map[string]*model.Subscription),
	}
}

func (m *MockSubscriptionService) Create(sub *model.Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()
	sub.ID = id
	m.data[id] = sub
	return nil
}

func (m *MockSubscriptionService) Get(id string) (*model.Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sub, ok := m.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return sub, nil
}

func (m *MockSubscriptionService) List(userID string) ([]model.Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var res []model.Subscription
	for _, sub := range m.data {
		if userID == "" || sub.UserID == userID {
			res = append(res, *sub)
		}
	}
	return res, nil
}

func (m *MockSubscriptionService) Update(id string, upd *model.UpdateSubscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sub, ok := m.data[id]
	if !ok {
		return errors.New("not found")
	}

	if upd.ServiceName != nil {
		sub.ServiceName = *upd.ServiceName
	}
	if upd.Price != nil {
		sub.Price = *upd.Price
	}
	if upd.UserID != nil {
		sub.UserID = *upd.UserID
	}
	if upd.StartDate != nil {
		sub.StartDate = *upd.StartDate
	}
	if upd.EndDate != nil {
		sub.EndDate = upd.EndDate
	}

	return nil
}

func (m *MockSubscriptionService) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[id]; !ok {
		return errors.New("not found")
	}
	delete(m.data, id)
	return nil
}

func (m *MockSubscriptionService) Aggregate(from, to time.Time, userID, service string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sum := 0
	for _, sub := range m.data {
		if sub.StartDate.Time.Before(from) || sub.StartDate.Time.After(to) {
			continue
		}
		if userID != "" && sub.UserID != userID {
			continue
		}
		if service != "" && sub.ServiceName != service {
			continue
		}
		sum += sub.Price
	}
	return sum, nil
}

var mockSvc *MockSubscriptionService

func setupTestRouter() *chi.Mux {
	mockSvc = NewMockSubscriptionService()
	logger := zap.NewExample().Sugar()
	r := router.NewRouter(time.Duration(30)*time.Second, mockSvc, logger)
	return r
}

func TestCreateSubscription(t *testing.T) {
	r := setupTestRouter()

	payload := map[string]interface{}{
		"service_name": "Netflix",
		"price":        500,
		"user_id":      "user-1",
		"start_date":   "07-2025",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/subs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Subscription
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Netflix", resp.ServiceName)
	assert.Equal(t, 500, resp.Price)
	assert.Equal(t, "user-1", resp.UserID)
	assert.NotEmpty(t, resp.ID)
}

func TestListSubscriptions(t *testing.T) {
	r := setupTestRouter()

	mockSvc.Create(&model.Subscription{
		ServiceName: "Spotify",
		Price:       300,
		UserID:      "user-2",
		StartDate:   model.MonthYear{Time: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/subs", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var subs []model.Subscription
	err := json.Unmarshal(w.Body.Bytes(), &subs)
	assert.NoError(t, err)
	assert.True(t, len(subs) >= 1)
}

func TestUpdateSubscription(t *testing.T) {
	r := setupTestRouter()

	sub := &model.Subscription{
		ServiceName: "Apple Music",
		Price:       400,
		UserID:      "user-3",
		StartDate:   model.MonthYear{Time: time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)},
	}
	mockSvc.Create(sub)

	updatePayload := map[string]interface{}{
		"price": 450,
	}
	body, _ := json.Marshal(updatePayload)

	url := "/api/v1/subs/" + sub.ID
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated model.Subscription
	err := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.NoError(t, err)
	assert.Equal(t, 450, updated.Price)
}

func TestDeleteSubscription(t *testing.T) {
	r := setupTestRouter()

	sub := &model.Subscription{
		ServiceName: "Deezer",
		Price:       200,
		UserID:      "user-4",
		StartDate:   model.MonthYear{Time: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)},
	}
	mockSvc.Create(sub)

	url := "/api/v1/subs/" + sub.ID
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err := mockSvc.Get(sub.ID)
	assert.Error(t, err)
}

func TestAggregateSubscription(t *testing.T) {
	r := setupTestRouter()

	mockSvc.Create(&model.Subscription{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      "user-5",
		StartDate:   model.MonthYear{Time: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
	})
	mockSvc.Create(&model.Subscription{
		ServiceName: "Netflix",
		Price:       200,
		UserID:      "user-5",
		StartDate:   model.MonthYear{Time: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)},
	})

	url := "/api/v1/subs/aggregate?from=2025-01&to=2025-02&user_id=user-5&service_name=Netflix"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 300, resp["total"])
}
