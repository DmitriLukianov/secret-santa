package v1

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"secret-santa-backend/internal/controller/http/middleware"
	"secret-santa-backend/internal/oauth"
)

func NewRouter(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	eventHandler *EventHandler,
	participantHandler *ParticipantHandler,
	assignmentHandler *AssignmentHandler,
	wishlistHandler *WishlistHandler,
	invitationHandler *InvitationHandler,
	chatHandler *ChatHandler,
	jwtManager *oauth.JWTManager,
	log *slog.Logger,
) *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.TimeoutMiddleware(10 * time.Second))

	router.Route("/auth", func(r chi.Router) {
		r.Get("/login", authHandler.Login)
		r.Get("/callback", authHandler.Callback)

		// === Новые эндпоинты для входа по email + OTP ===
		r.Post("/send-otp", authHandler.SendOTP)
		r.Post("/verify-otp", authHandler.VerifyOTP)
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(jwtManager, log).Handler)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.GetUsers)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
			r.Get("/me/events", userHandler.GetMyEvents)
		})

		r.Route("/events", func(r chi.Router) {
			r.Post("/", eventHandler.CreateEvent)
			r.Get("/", eventHandler.GetEvents)
			r.Get("/{id}", eventHandler.GetEventByID)
			r.Put("/{id}", eventHandler.UpdateEvent)
			r.Delete("/{id}", eventHandler.DeleteEvent)

			r.Post("/{id}/open-invitation", eventHandler.OpenInvitation)
			r.Post("/{id}/close-registration", eventHandler.CloseRegistration)
			r.Post("/{id}/start-drawing", eventHandler.StartDrawing)
			r.Post("/{id}/finish", eventHandler.FinishEvent)
			r.Post("/{id}/cancel", eventHandler.CancelEvent)

			r.Post("/{eventId}/participants", participantHandler.Add)
			r.Get("/{eventId}/participants", participantHandler.GetByEvent)
			r.Post("/{eventId}/assign", assignmentHandler.Draw)
			r.Get("/{eventId}/assignments", assignmentHandler.GetByEvent)

			r.Route("/{eventId}/chat", func(r chi.Router) {
				r.Get("/recipient", chatHandler.GetRecipientChat)
				r.Get("/sender", chatHandler.GetSenderChat)
				r.Post("/messages", chatHandler.SendMessage)
			})
		})

		r.Post("/participants/{id}/gift-sent", participantHandler.MarkGiftSent)
		r.Delete("/participants/{id}", participantHandler.Delete)

		r.Route("/users/me/wishlist", func(r chi.Router) {
			r.Post("/", wishlistHandler.Create)
			r.Get("/", wishlistHandler.GetByUser)
		})

		r.Route("/wishlists/{participantId}", func(r chi.Router) {
			r.Get("/", wishlistHandler.GetByParticipant)
		})

		r.Route("/wishlists/{wishlistId}/items", func(r chi.Router) {
			r.Post("/", wishlistHandler.AddItem)
			r.Get("/", wishlistHandler.GetItems)
			r.Put("/{itemId}", wishlistHandler.UpdateItem)
			r.Delete("/{itemId}", wishlistHandler.DeleteItem)
		})

		r.Route("/invitations", func(r chi.Router) {
			r.Post("/generate", invitationHandler.GenerateInvite)
		})
		r.Post("/invite/join", invitationHandler.JoinByInvite)
	})

	return router
}
