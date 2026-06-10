package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"frontdev333/gym/internal/config"
	"frontdev333/gym/internal/handler"
	"frontdev333/gym/internal/repository"
	"frontdev333/gym/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		slog.Error("db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		slog.Error("db ping", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)
	workoutRepo := repository.NewWorkoutRepository(db)
	statisticsRepo := repository.NewStatisticsRepository(db)

	userService := service.NewUserService(userRepo)
	exerciseService := service.NewExerciseService(exerciseRepo)
	workoutService := service.NewWorkoutService(workoutRepo, userRepo, exerciseRepo)
	statisticsService := service.NewStatisticsService(statisticsRepo)

	userHandler := handler.NewUserHandler(userService)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)
	workoutHandler := handler.NewWorkoutHandler(workoutService)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(30 * time.Second))

	mux.Get("/health", handler.Health)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetAll)
			r.Post("/", userHandler.Create)
			r.Get("/{user_id}", userHandler.GetByID)
			r.Put("/{user_id}", userHandler.Update)
			r.Delete("/{user_id}", userHandler.Delete)

			r.Route("/{user_id}/workouts", func(r chi.Router) {
				r.Get("/", workoutHandler.GetAll)
				r.Post("/", workoutHandler.Create)
			})

			r.Get("/{user_id}/statistics", statisticsHandler.GetByUserID)
		})

		r.Route("/exercises", func(r chi.Router) {
			r.Get("/", exerciseHandler.GetAll)
			r.Post("/", exerciseHandler.Create)
			r.Get("/{exercise_id}", exerciseHandler.GetByID)
			r.Put("/{exercise_id}", exerciseHandler.Update)
			r.Delete("/{exercise_id}", exerciseHandler.Delete)
		})
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("server started", "port", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown", "error", err)
	}
}
