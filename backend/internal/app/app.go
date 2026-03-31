package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"secret-santa-backend/internal/config"
	"secret-santa-backend/internal/logger"
	"secret-santa-backend/internal/middleware"

	"secret-santa-backend/internal/oauth"

	postgres "secret-santa-backend/internal/repository/postgres"
	assignmentrepo "secret-santa-backend/internal/repository/postgres/assignment"
	eventrepo "secret-santa-backend/internal/repository/postgres/event"
	invitationrepo "secret-santa-backend/internal/repository/postgres/invitation"
	participantrepo "secret-santa-backend/internal/repository/postgres/participant"
	userrepo "secret-santa-backend/internal/repository/postgres/user"
	wishlistrepo "secret-santa-backend/internal/repository/postgres/wishlist"

	assignmentusecase "secret-santa-backend/internal/usecase/assignment"
	authusecase "secret-santa-backend/internal/usecase/auth"
	eventusecase "secret-santa-backend/internal/usecase/event"
	invitationusecase "secret-santa-backend/internal/usecase/invitation"
	participantusecase "secret-santa-backend/internal/usecase/participant"
	userusecase "secret-santa-backend/internal/usecase/user"
	wishlistusecase "secret-santa-backend/internal/usecase/wishlist"

	v1 "secret-santa-backend/internal/controller/http/v1"
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	server *http.Server
}

func New() *App {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, cfg.AppEnv)

	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		panic(err)
	}

	// ==================== Репозитории ====================
	userRepo := userrepo.New(db)
	eventRepo := eventrepo.New(db)
	participantRepo := participantrepo.New(db)
	assignmentRepo := assignmentrepo.New(db)
	wishlistRepo := wishlistrepo.New(db)
	invitationRepo := invitationrepo.New(db)

	// ==================== UseCases ====================
	userUC := userusecase.NewWithLogger(userRepo, log)
	authUC := authusecase.NewWithLogger(userUC, log)

	eventUC := eventusecase.NewWithLogger(eventRepo, log)
	participantUC := participantusecase.NewWithLogger(participantRepo, log)

	assignmentUC := assignmentusecase.NewWithLogger(
		assignmentRepo,
		participantRepo,
		eventRepo,
		log,
	)

	// FIXED: теперь передаём participantRepo (нужен для GetForUser)
	wishlistUC := wishlistusecase.NewWithLogger(wishlistRepo, participantRepo, assignmentRepo, log)

	invitationUC := invitationusecase.NewWithLogger(
		invitationRepo,
		eventRepo,
		participantUC,
		log,
	)

	// ==================== Handlers ====================
	userHandler := v1.NewUserHandler(userUC, eventUC)
	eventHandler := v1.NewEventHandler(eventUC)
	participantHandler := v1.NewParticipantHandler(participantUC)
	wishlistHandler := v1.NewWishlistHandler(wishlistUC, participantUC)
	assignmentHandler := v1.NewAssignmentHandler(assignmentUC)
	invitationHandler := v1.NewInvitationHandler(invitationUC)

	// ==================== Auth ====================
	jwtManager, err := oauth.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL)
	if err != nil {
		log.Error("failed to create JWT manager", slog.String("error", err.Error()))
		panic(err)
	}

	authProvider, err := oauth.New(cfg)
	if err != nil {
		log.Error("failed to create oauth provider", slog.String("error", err.Error()))
		panic(err)
	}

	authHandler := v1.NewAuthHandler(authProvider, jwtManager, authUC)

	// ==================== Router ====================
	r := chi.NewRouter()

	// 🔥 Глобальные middleware (ВАЖЕН ПОРЯДОК)
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.TimeoutMiddleware(10 * time.Second))

	// ==================== Public routes ====================
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", authHandler.Login)
		r.Get("/callback", authHandler.Callback)
	})

	// ==================== Protected routes ====================
	r.Group(func(r chi.Router) {
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
		})

		r.Post("/participants/{id}/gift-sent", participantHandler.MarkGiftSent)
		r.Delete("/participants/{id}", participantHandler.Delete)

		r.Route("/users/{userId}/wishlist", func(r chi.Router) {
			r.Post("/", wishlistHandler.Create)
			r.Get("/", wishlistHandler.GetByUser)
		})

		r.Route("/wishlists/{wishlistId}/items", func(r chi.Router) {
			r.Post("/", wishlistHandler.AddItem)
			r.Get("/", wishlistHandler.GetItems)
		})

		r.Route("/invitations", func(r chi.Router) {
			r.Post("/generate", invitationHandler.GenerateInvite)
		})

		r.Post("/invite/join", invitationHandler.JoinByInvite)
	})

	return &App{
		cfg: cfg,
		log: log,
		server: &http.Server{
			Addr:    ":" + cfg.AppPort,
			Handler: r,
		},
	}
}

func (a *App) Run() error {
	a.log.Info("🚀 Server running on :" + a.cfg.AppPort)
	return a.server.ListenAndServe()
}
