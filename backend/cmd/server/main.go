package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	// repos
	"secret-santa-backend/internal/middleware"
	assignmentrepo "secret-santa-backend/internal/repository/postgres/assignment"
	eventrepo "secret-santa-backend/internal/repository/postgres/event"
	participantrepo "secret-santa-backend/internal/repository/postgres/participant"
	userrepo "secret-santa-backend/internal/repository/postgres/user"
	wishlistrepo "secret-santa-backend/internal/repository/postgres/wishlist"

	// db
	postgres "secret-santa-backend/internal/repository/postgres"

	// usecases
	assignmentusecase "secret-santa-backend/internal/usecase/assignment"
	eventusecase "secret-santa-backend/internal/usecase/event"
	participantusecase "secret-santa-backend/internal/usecase/participant"
	userusecase "secret-santa-backend/internal/usecase/user"
	wishlistusecase "secret-santa-backend/internal/usecase/wishlist"

	// 🔥 AUTH
	oauth "secret-santa-backend/internal/auth"
	authusecase "secret-santa-backend/internal/usecase/auth"

	// handlers
	v1 "secret-santa-backend/internal/controller/http/v1"
)

func main() {

	// ======================
	// ENV
	// ======================
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// ======================
	// DB
	// ======================
	db, err := postgres.NewDB(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ======================
	// OAUTH (VK)
	// ======================
	provider := oauth.NewGitHubProvider(
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		os.Getenv("GITHUB_REDIRECT_URL"),
	)

	// ======================
	// REPOSITORIES
	// ======================
	userRepo := userrepo.New(db)
	eventRepo := eventrepo.New(db)
	participantRepo := participantrepo.New(db)
	assignmentRepo := assignmentrepo.New(db)
	wishlistRepo := wishlistrepo.New(db)

	// ======================
	// USECASES
	// ======================
	userUC := userusecase.New(userRepo)
	eventUC := eventusecase.New(eventRepo)
	participantUC := participantusecase.New(participantRepo)
	assignmentUC := assignmentusecase.New(assignmentRepo, participantRepo)
	wishlistUC := wishlistusecase.New(wishlistRepo)

	// 🔥 AUTH USECASE
	authUC := authusecase.New(userRepo)

	// ======================
	// HANDLERS
	// ======================
	userHandler := v1.NewUserHandler(userUC)
	eventHandler := v1.NewEventHandler(eventUC)
	participantHandler := v1.NewParticipantHandler(participantUC)
	assignmentHandler := v1.NewAssignmentHandler(assignmentUC)
	wishlistHandler := v1.NewWishlistHandler(wishlistUC)

	// 🔥 AUTH HANDLER
	authHandler := v1.NewAuthHandler(provider, jwtManager, authUC)

	// ======================
	// ROUTER
	// ======================
	r := chi.NewRouter()

	// ======================
	// 🔐 AUTH (ПУБЛИЧНЫЕ)
	// ======================
	r.Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)

	// ======================
	// 🔒 PROTECTED ROUTES
	// ======================
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(jwtManager).Handler)

		// USERS
		r.Post("/users", userHandler.CreateUser)
		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)

		// EVENTS
		r.Post("/events", eventHandler.CreateEvent)
		r.Get("/events", eventHandler.GetEvents)
		r.Get("/events/{id}", eventHandler.GetEventByID)
		r.Put("/events/{id}", eventHandler.UpdateEvent)
		r.Delete("/events/{id}", eventHandler.DeleteEvent)

		// PARTICIPANTS
		r.Post("/events/{eventId}/participants", participantHandler.Add)
		r.Get("/events/{eventId}/participants", participantHandler.GetByEvent)
		r.Delete("/participants/{id}", participantHandler.Delete)

		// ASSIGNMENT
		r.Post("/events/{eventId}/assign", assignmentHandler.Draw)
		r.Get("/events/{eventId}/assignments", assignmentHandler.GetByEvent)

		// WISHLIST
		r.Post("/users/{userId}/wishlist", wishlistHandler.Create)
		r.Get("/users/{userId}/wishlist", wishlistHandler.GetByUser)
		r.Delete("/wishlist/{id}", wishlistHandler.Delete)
	})

	// ======================
	// START SERVER
	// ======================
	log.Println("🚀 Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
