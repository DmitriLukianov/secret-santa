package main

import (
	"log"
	"net/http"
	"os"
	"secret-santa-backend/internal/handlers"
	"secret-santa-backend/internal/repository"
	"secret-santa-backend/internal/repository/postgres"
	"secret-santa-backend/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := repository.NewDB(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	eventRepo := postgres.NewEventRepository(db)
	participantRepo := postgres.NewParticipantRepository(db)
	wishlistRepo := postgres.NewWishlistRepository(db)
	assignmentRepo := postgres.NewAssignmentRepository(db)

	userService := services.NewUserService(userRepo)
	eventService := services.NewEventService(eventRepo)
	participantService := services.NewParticipantService(participantRepo)
	wishlistService := services.NewWishlistService(wishlistRepo)
	assignmentService := services.NewAssignmentService(assignmentRepo)

	userHandler := handlers.NewUserHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService)
	participantHandler := handlers.NewParticipantHandler(participantService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	assignmentHandler := handlers.NewAssignmentHandler(assignmentService)

	r := chi.NewRouter()

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

	r.Post("/events/{event_id}/participants", participantHandler.JoinEvent)
	r.Get("/events/{event_id}/participants", participantHandler.GetParticipants)
	r.Delete("/events/{event_id}/participants/{user_id}", participantHandler.LeaveEvent)

	r.Post("/events/{event_id}/wishlist", wishlistHandler.CreateWishlist)
	r.Get("/events/{event_id}/wishlist/{user_id}", wishlistHandler.GetWishlist)
	r.Put("/wishlist/{id}", wishlistHandler.UpdateWishlist)

	r.Post("/assignments", assignmentHandler.Create)
	r.Get("/assignments/my", assignmentHandler.GetMyAssignment)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
