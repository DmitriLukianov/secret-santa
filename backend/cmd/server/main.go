package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	authpkg "secret-santa-backend/internal/auth"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	jwtTTL, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil || jwtTTL == 0 {
		jwtTTL = 24 * time.Hour
	}

	jwtManager, err := authpkg.NewJWTManager(jwtSecret, jwtTTL)
	if err != nil {
		log.Fatal("failed to create JWT manager:", err)
	}

	provider := authpkg.NewGitHubProvider(
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		os.Getenv("GITHUB_REDIRECT_URL"),
	)

	// ==================== REPOSITORIES ====================
	userRepo := userrepo.New(db)
	eventRepo := eventrepo.New(db)
	participantRepo := participantrepo.New(db)
	assignmentRepo := assignmentrepo.New(db)
	wishlistRepo := wishlistrepo.New(db)

	// ==================== USECASES ====================
	userUC := userusecase.New(userRepo)
	authUC := authusecase.New(userUC)

	eventUC := eventusecase.New(eventRepo)
	participantUC := participantusecase.New(participantRepo)
	assignmentUC := assignmentusecase.New(assignmentRepo, participantRepo)
	wishlistUC := wishlistusecase.New(wishlistRepo)

	// ==================== HANDLERS ====================
	userHandler := v1.NewUserHandler(userUC)
	eventHandler := v1.NewEventHandler(eventUC)
	participantHandler := v1.NewParticipantHandler(participantUC)
	assignmentHandler := v1.NewAssignmentHandler(assignmentUC)
	wishlistHandler := v1.NewWishlistHandler(wishlistUC)
	authHandler := v1.NewAuthHandler(provider, jwtManager, authUC)

	r := chi.NewRouter()

	// Public routes
	r.Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(jwtManager).Handler)

		r.Post("/users", userHandler.CreateUser)
		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)

		r.Post("/events", eventHandler.CreateEvent)
		r.Get("/events", eventHandler.GetEvents)
		r.Get("/events/{id}", eventHandler.GetEventByID)
		r.Put("/events/{id}", eventHandler.UpdateEvent)
		r.Delete("/events/{id}", eventHandler.DeleteEvent)

		r.Post("/events/{eventId}/participants", participantHandler.Add)
		r.Get("/events/{eventId}/participants", participantHandler.GetByEvent)
		r.Delete("/participants/{id}", participantHandler.Delete)

		r.Post("/events/{eventId}/assign", assignmentHandler.Draw)
		r.Get("/events/{eventId}/assignments", assignmentHandler.GetByEvent)

		r.Post("/users/{userId}/wishlist", wishlistHandler.Create)
		r.Get("/users/{userId}/wishlist", wishlistHandler.GetByUser)
		r.Delete("/wishlist/{id}", wishlistHandler.Delete)
	})

	log.Println("🚀 Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
