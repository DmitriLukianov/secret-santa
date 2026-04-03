package app

import (
	"log/slog"
	"net/http"

	"secret-santa-backend/internal/config"
	"secret-santa-backend/internal/logger"

	"secret-santa-backend/internal/oauth"

	"secret-santa-backend/internal/database"

	assignmentrepo "secret-santa-backend/internal/repository/postgres/assignment"
	chatrepo "secret-santa-backend/internal/repository/postgres/chat"
	eventrepo "secret-santa-backend/internal/repository/postgres/event"
	invitationrepo "secret-santa-backend/internal/repository/postgres/invitation"
	participantrepo "secret-santa-backend/internal/repository/postgres/participant"
	userrepo "secret-santa-backend/internal/repository/postgres/user"
	wishlistrepo "secret-santa-backend/internal/repository/postgres/wishlist"

	assignmentusecase "secret-santa-backend/internal/usecase/assignment"
	authusecase "secret-santa-backend/internal/usecase/auth"
	chatusecase "secret-santa-backend/internal/usecase/chat"
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

	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		panic(err)
	}

	userRepo := userrepo.New(db)
	eventRepo := eventrepo.New(db)
	participantRepo := participantrepo.New(db)
	assignmentRepo := assignmentrepo.New(db)
	wishlistRepo := wishlistrepo.New(db)
	invitationRepo := invitationrepo.New(db)
	chatRepo := chatrepo.New(db)

	userUC := userusecase.NewWithLogger(userRepo, log)
	authUC := authusecase.NewWithLogger(userUC, log)

	eventUC := eventusecase.NewWithLogger(eventRepo, participantRepo, log)
	participantUC := participantusecase.NewWithLogger(participantRepo, log)

	assignmentUC := assignmentusecase.NewWithLogger(assignmentRepo, participantRepo, eventRepo, log)
	wishlistUC := wishlistusecase.NewWithLogger(wishlistRepo, participantRepo, assignmentRepo, log)

	invitationUC := invitationusecase.NewWithLogger(invitationRepo, eventRepo, participantUC, log)
	chatUC := chatusecase.NewWithLogger(chatRepo, participantRepo, assignmentRepo, log)

	userHandler := v1.NewUserHandler(userUC, eventUC)
	eventHandler := v1.NewEventHandler(eventUC)
	participantHandler := v1.NewParticipantHandler(participantUC)
	wishlistHandler := v1.NewWishlistHandler(wishlistUC, participantUC)
	assignmentHandler := v1.NewAssignmentHandler(assignmentUC)
	invitationHandler := v1.NewInvitationHandler(invitationUC)
	chatHandler := v1.NewChatHandler(chatUC)

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

	router := v1.NewRouter(
		authHandler,
		userHandler,
		eventHandler,
		participantHandler,
		assignmentHandler,
		wishlistHandler,
		invitationHandler,
		chatHandler,
		jwtManager,
		log,
	)

	return &App{
		cfg: cfg,
		log: log,
		server: &http.Server{
			Addr:    ":" + cfg.AppPort,
			Handler: router,
		},
	}
}

func (a *App) Run() error {
	a.log.Info("🚀 Server running on :" + a.cfg.AppPort)
	return a.server.ListenAndServe()
}
