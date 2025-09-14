package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nearbygems/subscription-service/internal/model"
	"github.com/sirupsen/logrus"
)

type mockStore struct {
	items map[uuid.UUID]model.Subscription
}

var ErrNotFound = errors.New("not found")

func newMockStore() *mockStore {
	return &mockStore{items: make(map[uuid.UUID]model.Subscription)}
}

func (m *mockStore) Create(sub *model.Subscription) error {
	m.items[sub.ID] = *sub
	return nil
}

func (m *mockStore) Get(id uuid.UUID) (*model.Subscription, error) {
	sub, ok := m.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &sub, nil
}

func (m *mockStore) Update(sub *model.Subscription) error {
	m.items[sub.ID] = *sub
	return nil
}

func (m *mockStore) Delete(id uuid.UUID) error {
	delete(m.items, id)
	return nil
}

func (m *mockStore) List(limit, offset int, userID *uuid.UUID, serviceName *string) ([]model.Subscription, error) {
	result := []model.Subscription{}
	for _, sub := range m.items {
		result = append(result, sub)
	}
	return result, nil
}

func (m *mockStore) Summary(periodFrom, periodTo string, userID *uuid.UUID, serviceName *string) (int, error) {
	return 1000, nil
}

// -------------------- Тесты --------------------

func TestCreateSubscription(t *testing.T) {
	store := newMockStore()
	h := NewHandler(store, logrus.New())

	sub := model.Subscription{
		ID:          uuid.New(),
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   "07-2025",
	}

	body, _ := json.Marshal(sub)
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("ожидали 201, получили %d", w.Code)
	}
}

func TestGetNotFound(t *testing.T) {
	store := newMockStore()
	h := NewHandler(store, logrus.New())

	missingID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+missingID.String(), nil)
	w := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", missingID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	h.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("ожидали 404, получили %d", w.Code)
	}
}

func TestSummary(t *testing.T) {
	store := newMockStore()
	h := NewHandler(store, logrus.New())

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/summary?period_from=01-2025&period_to=12-2025", nil)
	w := httptest.NewRecorder()

	h.Summary(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", w.Code)
	}
}
