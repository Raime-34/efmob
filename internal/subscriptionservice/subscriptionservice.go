package subscriptionservice

import (
	"context"
	"efmob/internal/cfg"
	"efmob/internal/dto"
	"efmob/internal/migrate"
	"efmob/internal/repositories"
	"efmob/logger"
	"fmt"
	"net/http"

	_ "efmob/docs" // swagger generated docs

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type SubService struct {
	repo      Repositories
	validator *validator.Validate
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
		repo:      repositories.NewRepo(connPool),
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
	service.mountHandlers()

	return service
}

func (s *SubService) mountHandlers() {
	r := chi.NewRouter()
	r.Use(requestLogMiddleware)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscription", func(r chi.Router) {
			r.Get("/sum", s.SubscriptionFilteredSumHandler)
			r.Get("/", s.ListSubscriptionsHandler)
			r.Post("/", s.InsertSubscriptionHandler)
			r.Delete("/", s.DeleteSubscriptionHandler)
			r.Put("/", s.UpdateSubscriptionHandler)
			r.Patch("/", s.PatchSubscriptionHandler)
		})
	})

	port := cfg.GetConfig().Port
	go http.ListenAndServe(fmt.Sprintf(":%v", port), r)
	fmt.Printf("Сваггер доступен по адресу: http://localhost:%v/swagger/\n", port)
}

//go:generate mockgen -source=subscriptionservice.go -destination=../../mocks/repositories.go -package=mock
type Repositories interface {
	InsertSubscriptionInfo(context.Context, dto.CreateOrUpdateSubscriptionRequest) error
	DeleteSubscriptionInfo(context.Context, dto.DeleteSubscriptionRequest) error
	UpdateSubscriptionInfo(context.Context, dto.CreateOrUpdateSubscriptionRequest) error
	PatchSubscriptionInfo(context.Context, dto.PatchSubscriptionRequest) error
	ListSubscriptionsByUserID(context.Context, string, int, int) ([]dto.SubscriptionListItem, int, error)
	SumSubscriptionPriceByFilter(context.Context, string, string, string, string) (int64, error)
}
