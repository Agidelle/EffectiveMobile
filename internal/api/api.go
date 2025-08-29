package api

import (
	"context"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	service SubService
}

type SubService interface {
	CloseDB()
	ListSubscriptions(ctx context.Context) ([]*domain.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int) (*domain.Subscription, error)
	CreateSubscription(ctx context.Context, input *domain.SubscriptionInput) (int64, error)
	UpdateSubscription(ctx context.Context, id int, input *domain.SubscriptionInput) error
	DeleteSubscription(ctx context.Context, id int) error
	GetSubscriptionsSummary(ctx context.Context) (interface{}, error)
}

func NewHandler(s SubService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes(r chi.Router) {
	// Роуты для управления подписками
	r.Get("/api/subscriptions", h.ListSubscriptions)          // список подписок
	r.Get("/api/subscriptions/{id}", h.GetSubscriptionByID)   // получение подписки по ID
	r.Post("/api/subscriptions", h.CreateSubscription)        // создание новой подписки
	r.Put("/api/subscriptions/{id}", h.UpdateSubscription)    // обновление подписки по ID
	r.Delete("/api/subscriptions/{id}", h.DeleteSubscription) // удаление подписки по ID

	r.Get("/api/subscriptions/summary", h.GetSubscriptionsSummary) // сводная информация по подпискам
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetSubscriptionsSummary(w http.ResponseWriter, r *http.Request) {
}
