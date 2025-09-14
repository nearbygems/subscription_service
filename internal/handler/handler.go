package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nearbygems/subscription-service/internal/model"
	"github.com/nearbygems/subscription-service/internal/store"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	store store.Store
	log   *logrus.Logger
}

func NewHandler(s store.Store, log *logrus.Logger) *Handler {
	return &Handler{store: s, log: log}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	sub := model.Subscription{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub.ID = uuid.New()
	if err := h.store.Create(&sub); err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sub)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	sub, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(sub)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	sub := model.Subscription{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub.ID = id
	if err := h.store.Update(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(sub)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	var userID *uuid.UUID
	if u := r.URL.Query().Get("user_id"); u != "" {
		id, err := uuid.Parse(u)
		if err == nil {
			userID = &id
		}
	}
	var serviceName *string
	if s := r.URL.Query().Get("service_name"); s != "" {
		serviceName = &s
	}
	list, err := h.store.List(limit, offset, userID, serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	periodFrom := r.URL.Query().Get("period_from")
	periodTo := r.URL.Query().Get("period_to")
	if periodFrom == "" || periodTo == "" {
		http.Error(w, "period_from and period_to required",
			http.StatusBadRequest)
		return
	}
	var userID *uuid.UUID
	if u := r.URL.Query().Get("user_id"); u != "" {
		id, err := uuid.Parse(u)
		if err == nil {
			userID = &id
		}
	}
	var serviceName *string
	if s := r.URL.Query().Get("service_name"); s != "" {
		serviceName = &s
	}
	total, err := h.store.Summary(periodFrom, periodTo, userID, serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]int{"total": total})
}
