package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	authpkg "secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/logger"
	"secret-santa-backend/internal/middleware"

	postgres "secret-santa-backend/internal/repository/postgres"
	assignmentrepo "secret-santa-backend/internal/repository/postgres/assignment"
	eventrepo "secret-santa-backend/internal/repository/postgres/event"
	participantrepo "secret-santa-backend/internal/repository/postgres/participant"
	userrepo "secret-santa-backend/internal/repository/postgres/user"
	wishlistrepo "secret-santa-backend/internal/repository/postgres/wishlist"

	assignmentusecase "secret-santa-backend/internal/usecase/assignment"
	authusecase "secret-santa-backend/internal/usecase/auth"
	eventusecase "secret-santa-backend/internal/usecase/event"
	participantusecase "secret-santa-backend/internal/usecase/participant"
	userusecase "secret-santa-backend/internal/usecase/user"
	wishlistusecase "secret-santa-backend/internal/usecase/wishlist"

	v1 "secret-santa-backend/internal/controller/http/v1"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := postgres.NewDB(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	stage := os.Getenv("APP_ENV")
	if stage == "" {
		stage = "local"
	}

	log := logger.New(logLevel, stage)
	log.Info("logger initialized",
		slog.String("level", logLevel),
		slog.String("stage", stage),
	)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Error("JWT_SECRET is not set")
		os.Exit(1)
	}

	jwtTTL, _ := time.ParseDuration(os.Getenv("JWT_TTL"))
	if jwtTTL == 0 {
		jwtTTL = 24 * time.Hour
	}

	jwtManager, err := authpkg.NewJWTManager(jwtSecret, jwtTTL)
	if err != nil {
		log.Error("failed to create JWT manager", slog.String("error", err.Error()))
		os.Exit(1)
	}

	provider := authpkg.NewGitHubProvider(
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		os.Getenv("GITHUB_REDIRECT_URL"),
	)

	userRepo := userrepo.New(db)
	eventRepo := eventrepo.New(db)
	participantRepo := participantrepo.New(db)
	assignmentRepo := assignmentrepo.New(db)
	wishlistRepo := wishlistrepo.New(db)

	userUC := userusecase.NewWithLogger(userRepo, log)
	authUC := authusecase.New(userUC)

	eventUC := eventusecase.NewWithLogger(eventRepo, log)
	participantUC := participantusecase.NewWithLogger(participantRepo, log)

	assignmentUC := assignmentusecase.NewWithLogger(assignmentRepo, participantRepo, eventUC, log)

	wishlistUC := wishlistusecase.NewWithLogger(wishlistRepo, assignmentUC, log)

	userHandler := v1.NewUserHandler(userUC, eventUC)
	eventHandler := v1.NewEventHandler(eventUC)
	participantHandler := v1.NewParticipantHandler(participantUC)
	wishlistHandler := v1.NewWishlistHandler(wishlistUC, participantUC)
	assignmentHandler := v1.NewAssignmentHandler(assignmentUC)
	authHandler := v1.NewAuthHandler(provider, jwtManager, authUC)

	r := chi.NewRouter()

	r.Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)

	r.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(jwtManager).Handler)

		r.Post("/users", userHandler.CreateUser)
		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)

		r.Get("/me/events", userHandler.GetMyEvents)

		r.Post("/events", eventHandler.CreateEvent)
		r.Get("/events", eventHandler.GetEvents)
		r.Get("/events/{id}", eventHandler.GetEventByID)
		r.Put("/events/{id}", eventHandler.UpdateEvent)
		r.Delete("/events/{id}", eventHandler.DeleteEvent)
		r.Post("/events/{id}/finish", eventHandler.FinishEvent)

		r.Post("/events/{eventId}/participants", participantHandler.Add)
		r.Get("/events/{eventId}/participants", participantHandler.GetByEvent)
		r.Post("/participants/{id}/gift-sent", participantHandler.MarkGiftSent)
		r.Delete("/participants/{id}", participantHandler.Delete)

		r.Post("/users/{userId}/wishlist", wishlistHandler.Create)
		r.Get("/users/{userId}/wishlist", wishlistHandler.GetByUser)
		r.Post("/wishlists/{wishlistId}/items", wishlistHandler.AddItem)
		r.Get("/wishlists/{wishlistId}/items", wishlistHandler.GetItems)

		r.Post("/events/{eventId}/assign", assignmentHandler.Draw)
		r.Get("/events/{eventId}/assignments", assignmentHandler.GetByEvent)
	})

	log.Info("🚀 Server running on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
