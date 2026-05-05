package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"

	"effective-task/internal/config"
	"effective-task/internal/handler"
	mg "effective-task/internal/migrations"
	"effective-task/internal/repository"
	"effective-task/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leosunmo/zapchi"
	"go.uber.org/zap"
)

const (
	migrationsDir = "migrations"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	connDBstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	mgr, err := mg.NewMigrator(migrationsFS, migrationsDir)
	if err != nil {
		logger.Fatal("failed to create migrator", zap.Error(err))
	}
	db, err := sql.Open("postgres", connDBstr)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	err = mgr.ApplyMigrations(db)
	if err != nil {
		logger.Fatal("failed to apply migrations", zap.Error(err))
	}
	db.Close()

	repo, err := repository.NewRepository(connDBstr)
	if err != nil {
		logger.Fatal("failed to create repository", zap.Error(err))
	}

	serv := service.NewService(repo)
	handlerInstance := handler.NewHandler(logger, serv)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(zapchi.Logger(logger, "router"))
	router.Use(corsMiddleware)
	router.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	router.Post("/create/", handlerInstance.Create)
	router.Post("/read/", handlerInstance.ReadBySubscriptionID)
	router.Delete("/delete/", handlerInstance.Delete)
	router.Post("/update/", handlerInstance.Update)
	router.Post("/list/", handlerInstance.ListByUserID)
	router.Post("/count-by-user-id/", handlerInstance.CountByUserID)
	router.Post("/count-by-service-name/", handlerInstance.CountByServiceName)
	logger.Info("server starting", zap.String("address", ":8080"))
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}
