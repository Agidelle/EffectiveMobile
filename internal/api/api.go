package api

import (
	"context"
	"encoding/json"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const tOutnormal = 3 * time.Second
const tOutlong = 10 * time.Second

type Handler struct {
	service SubService
}

type SubService interface {
	CloseDB()
	Search(ctx context.Context, filter *domain.Filter) ([]*domain.Subscription, error)
	CreateSubscription(ctx context.Context, input *domain.Subscription) error
	UpdateSubscription(ctx context.Context, input *domain.Subscription) error
	DeleteSubscription(ctx context.Context, filter *domain.Filter) error
	GetSubscriptionsSummary(ctx context.Context, filter *domain.Filter) (int, error)
}

func NewHandler(s SubService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes(r chi.Router) {
	// Роуты для управления подписками
	r.Get("/api/subscriptions", h.SearchSubscriptions)   // список подписок
	r.Post("/api/subscriptions", h.CreateSubscription)   // создание новой подписки
	r.Put("/api/subscriptions", h.UpdateSubscription)    // обновление подписки по ID
	r.Delete("/api/subscriptions", h.DeleteSubscription) // удаление подписки по ID

	r.Post("/api/subscriptions/summary", h.GetSubscriptionsSummary) // сводная информация по подпискам
}

// @Summary      Получить список подписок
// @Description  Возвращает список подписок с фильтрами и пагинацией
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        user_id      query     string  false  "ID пользователя"
// @Param        service_name query     string  false  "Название сервиса"
// @Param        price        query     int     false  "Цена"
// @Param        start_date   query     string  false  "Дата начала"
// @Param        end_date     query     string  false  "Дата конца"
// @Param        limit        query     int     false  "Лимит"
// @Param        offset       query     int     false  "Смещение"
// @Success      200  {array}  domain.Subscription
// @Failure      500  {string}  string  "internal server error"
// @Router       /api/subscriptions [get]
func (h *Handler) SearchSubscriptions(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter

	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		filter.UserID = &userID
	}
	serviceName := r.URL.Query().Get("service_name")
	if serviceName != "" {
		filter.ServiceName = &serviceName
	}
	priceStr := r.URL.Query().Get("price")
	if priceStr != "" {
		if price, err := strconv.Atoi(priceStr); err == nil {
			filter.Price = &price
		}
	}
	startDate := r.URL.Query().Get("start_date")
	if startDate != "" {
		filter.StartDate = &startDate
	}
	endDate := r.URL.Query().Get("end_date")
	if endDate != "" {
		filter.EndDate = &endDate
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = &limit
		}
	}
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = &offset
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), tOutlong)
	defer cancel()
	subs, err := h.service.Search(ctx, &filter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if subs == nil {
		subs = []*domain.Subscription{}
	}
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		slog.Error("Failed to encode to JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// @Summary      Создать подписку
// @Description  Создаёт новую подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      domain.SubscriptionInput  true  "Данные подписки"
// @Success      201  {string}  string  "created"
// @Failure      400  {string}  string  "error decode"
// @Failure      500  {string}  string  "internal error"
// @Router       /api/subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input domain.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "error decode", http.StatusBadRequest)
		return
	}

	opts := input.SubscriptionToOptions()
	sub := domain.NewSubscription(opts...)

	ctx, cancel := context.WithTimeout(r.Context(), tOutnormal)
	defer cancel()

	if err := h.service.CreateSubscription(ctx, sub); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// @Summary      Обновить подписку
// @Description  Обновляет существующую подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      domain.SubscriptionInput  true  "Данные подписки"
// @Success      200  {string}  string  "ok"
// @Failure      400  {string}  string  "error decode"
// @Failure      500  {string}  string  "internal error"
// @Router       /api/subscriptions [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var input domain.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "error decode", http.StatusBadRequest)
		return
	}

	opts := input.SubscriptionToOptions()
	sub := domain.NewSubscription(opts...)

	ctx, cancel := context.WithTimeout(r.Context(), tOutnormal)
	defer cancel()

	if err := h.service.UpdateSubscription(ctx, sub); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary      Удалить подписку
// @Description  Удаляет подписку по user_id и service_name
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        user_id      query     string  true  "ID пользователя"
// @Param        service_name query     string  true  "Название сервиса"
// @Success      200  {string}  string  "ok"
// @Failure      400  {string}  string  "user_id and service_name are required"
// @Failure      404  {string}  string  "subscription not found"
// @Failure      500  {string}  string  "internal server error"
// @Router       /api/subscriptions [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter

	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		filter.UserID = &userID
	}
	serviceName := r.URL.Query().Get("service_name")
	if serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if filter.UserID == nil || filter.ServiceName == nil {
		http.Error(w, "user_id and service_name are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), tOutnormal)
	defer cancel()

	err := h.service.DeleteSubscription(ctx, &filter)
	if err != nil {
		if err.Error() == "subscription not found" {
			http.Error(w, "subscription not found", http.StatusNotFound)
		} else {
			slog.Error("Failed to delete subscription", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary      Сводная информация по подпискам
// @Description  Возвращает сумму подписок по фильтру
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        filter  body  domain.Filter  true  "Фильтр"
// @Success      200  {object}  map[string]int
// @Failure      400  {string}  string  "invalid filter"
// @Failure      500  {string}  string  "internal server error"
// @Router       /api/subscriptions/summary [post]
func (h *Handler) GetSubscriptionsSummary(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		slog.Error("Failed to decode json", "error", err)
		http.Error(w, "invalid filter", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), tOutlong)
	defer cancel()
	total, err := h.service.GetSubscriptionsSummary(ctx, &filter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	resp := map[string]int{"total_price": total}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode to JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
