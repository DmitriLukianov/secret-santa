package v1

import (
	"log/slog"
	"net/http"
	"time"

	"secret-santa-backend/internal/controller/http/middleware"
	"secret-santa-backend/internal/oauth"

	"github.com/go-chi/chi/v5"
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

	// Публичные маршруты
	router.Route("/auth", func(r chi.Router) {
		r.Get("/login", authHandler.Login)
		r.Get("/callback", authHandler.Callback)
		r.Post("/send-otp", authHandler.SendOTP)
		r.Post("/verify-otp", authHandler.VerifyOTP)
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Все защищённые маршруты
	router.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(jwtManager, log).Handler)

		// --- Users ---
		// ВАЖНО: /me/* должны быть зарегистрированы ДО /{id},
		// иначе chi будет матчить "me" как UUID и возвращать 400.
		r.Route("/users", func(r chi.Router) {
			r.Get("/me", userHandler.GetMe)
			r.Patch("/me", userHandler.UpdateMe)
			r.Get("/me/events", userHandler.GetMyEvents)

			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.GetUsers)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// --- Events ---
		r.Route("/events", func(r chi.Router) {
			r.Post("/", eventHandler.CreateEvent)
			r.Get("/", eventHandler.GetEvents)
			r.Get("/{id}", eventHandler.GetEventByID)
			r.Put("/{id}", eventHandler.UpdateEvent)
			r.Delete("/{id}", eventHandler.DeleteEvent)

			// Смена статуса — только организатор (проверяется в usecase)
			r.Post("/{id}/open-invitation", eventHandler.OpenInvitation)
			r.Post("/{id}/close-registration", eventHandler.CloseRegistration)
			r.Post("/{id}/start-drawing", eventHandler.StartDrawing)
			r.Post("/{id}/finish", eventHandler.FinishEvent)
			r.Post("/{id}/cancel", eventHandler.CancelEvent)

			// Участники события
			r.Post("/{eventId}/participants", participantHandler.Add)
			r.Get("/{eventId}/participants", participantHandler.GetByEvent)

			// Жеребьёвка
			r.Post("/{eventId}/assign", assignmentHandler.Draw)
			r.Get("/{eventId}/assignments", assignmentHandler.GetByEvent)

			// Чат (анонимные вопросы тайному санте)
			r.Route("/{eventId}/chat", func(r chi.Router) {
				r.Get("/recipient", chatHandler.GetRecipientChat)
				r.Get("/sender", chatHandler.GetSenderChat)
				r.Post("/messages", chatHandler.SendMessage)
			})
		})

		// --- Participants (отдельные операции) ---
		r.Post("/participants/{id}/gift-sent", participantHandler.MarkGiftSent)
		r.Delete("/participants/{id}", participantHandler.Delete)

		// --- Wishlists ---
		// Вишлист текущего пользователя для события
		r.Route("/users/me/wishlist", func(r chi.Router) {
			r.Post("/", wishlistHandler.Create)
			r.Get("/", wishlistHandler.GetByUser)
		})

		// Вишлист конкретного участника (с проверкой видимости)
		r.Route("/wishlists/{participantId}", func(r chi.Router) {
			r.Get("/", wishlistHandler.GetByParticipant)
		})

		// Элементы вишлиста
		r.Route("/wishlists/{wishlistId}/items", func(r chi.Router) {
			r.Post("/", wishlistHandler.AddItem)
			r.Get("/", wishlistHandler.GetItems)
			r.Put("/{itemId}", wishlistHandler.UpdateItem)
			r.Delete("/{itemId}", wishlistHandler.DeleteItem)
		})

		// --- Invitations ---
		r.Route("/invitations", func(r chi.Router) {
			r.Post("/generate", invitationHandler.GenerateInvite)
		})
		r.Post("/invite/join", invitationHandler.JoinByInvite)
	})

	return router
}
