package subscriptionservice

import (
	"context"
	"efmob/internal/cfg"
	"efmob/internal/dto"
	"efmob/internal/migrate"
	"efmob/internal/repositories"
	"efmob/logger"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SubService struct {
	repo Repositories
}

func InitService() *SubService {
	dbConfig, err := pgxpool.ParseConfig(cfg.GetConfig().DbDSN)
	if err != nil {
		logger.Log().Fatal("parse db dsn", zap.Error(err))
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Log().Fatal("pgx pool", zap.Error(err))
	}

	migrate.MakeMigration(connPool)

	service := &SubService{
		repo: repositories.NewRepo(connPool),
	}
	service.mountHandlers()

	return service
}

func (s *SubService) mountHandlers() {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscription", func(r chi.Router) {
			r.Post("/", s.InsertSubscriptionHandler)
			r.Delete("/", s.DeleteSubscriptionHandler)
			r.Put("/", s.UpdateSubscriptionHandler)
			r.Patch("/", nil)
		})
	})

	go http.ListenAndServe(":8080", r)
}

type Repositories interface {
	InsertSubscriptionInfo(context.Context, dto.CreateSubscriptionDTO) error
	DeleteSubscriptionInfo(context.Context, dto.CreateSubscriptionDTO) error
	UpdateSubscriptionInfo(context.Context, dto.CreateSubscriptionDTO) error
}

func (s *SubService) InsertSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var subInfo dto.CreateSubscriptionDTO
	err := decoder.Decode(&subInfo)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler - Failed to parse request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.InsertSubscriptionInfo(r.Context(), subInfo)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *SubService) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var subInfo dto.CreateSubscriptionDTO
	err := decoder.Decode(&subInfo)
	if err != nil {
		logger.Log().Error("DeleteSubscriptionHandler - Failed to parse request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.DeleteSubscriptionInfo(r.Context(), subInfo)
	if err != nil {
		logger.Log().Error("DeleteSubscriptionHandler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *SubService) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var subInfo dto.CreateSubscriptionDTO
	err := decoder.Decode(&subInfo)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler - Failed to parse request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.UpdateSubscriptionInfo(r.Context(), subInfo)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
