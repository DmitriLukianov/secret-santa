package usecase

import (
	"context"
	"time"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UserUseCase interface {
	Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EventUseCase interface {
	Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	GetAll(ctx context.Context) ([]entity.Event, error)

	Update(ctx context.Context, id, userID uuid.UUID, input dto.UpdateEventInput) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status definitions.EventStatus) error
	OpenInvitation(ctx context.Context, id, userID uuid.UUID) error
	CloseRegistration(ctx context.Context, id, userID uuid.UUID) error
	StartDrawing(ctx context.Context, id, userID uuid.UUID) error
	Activate(ctx context.Context, id, userID uuid.UUID) error
	Finish(ctx context.Context, id, userID uuid.UUID) error
	Cancel(ctx context.Context, id, userID uuid.UUID) error
	GetMyEvents(ctx context.Context, userID uuid.UUID) ([]entity.Event, error)
}

type ParticipantUseCase interface {
	Create(ctx context.Context, eventID, userID uuid.UUID, role string) (entity.Participant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)

	MarkGiftSent(ctx context.Context, participantID, requesterID uuid.UUID) error

	Delete(ctx context.Context, id, requesterID uuid.UUID) error
	GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*entity.Participant, error)
}

type WishlistUseCase interface {
	Create(ctx context.Context, participantID uuid.UUID, visibility string) (entity.Wishlist, error)
	AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string, price *float64) (entity.WishlistItem, error)
	GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error)
	GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error)
	GetForUser(ctx context.Context, eventID, participantID, requesterID uuid.UUID) (*entity.Wishlist, error)
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.WishlistItem, error)
	UpdateItem(ctx context.Context, itemID uuid.UUID, title string, link, imageURL, comment *string, price *float64) (entity.WishlistItem, error)
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
	GetByID(ctx context.Context, wishlistID uuid.UUID) (*entity.Wishlist, error)
}

type AssignmentUseCase interface {
	Draw(ctx context.Context, eventID, userID uuid.UUID) error
	GetByEvent(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Assignment, error)
}

type InvitationUseCase interface {
	GenerateInvite(ctx context.Context, input dto.CreateInvitationInput, organizerID uuid.UUID) (dto.InvitationResponse, error)
	JoinByInvite(ctx context.Context, input dto.JoinByInvitationInput) error
	SendEmailInvitation(ctx context.Context, input dto.CreateInvitationInput, organizerID uuid.UUID, recipientEmail string) (dto.InvitationResponse, error)
}

type NotificationUseCase interface {
	Notify(ctx context.Context, userID uuid.UUID, notifType string, payload map[string]string) error
	GetForUser(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error)
	MarkAsRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}

type FriendshipUseCase interface {
	SendRequest(ctx context.Context, requesterID, addresseeID uuid.UUID) (entity.Friendship, error)
	AcceptRequest(ctx context.Context, friendshipID, userID uuid.UUID) error
	DeclineRequest(ctx context.Context, friendshipID, userID uuid.UUID) error
	RemoveFriend(ctx context.Context, friendshipID, userID uuid.UUID) error
	GetFriends(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error)
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error)
	AreFriends(ctx context.Context, userA, userB uuid.UUID) (bool, error)
}

type ChatUseCase interface {
	GetRecipientChat(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Message, error)
	GetSenderChat(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Message, error)
	SendMessage(ctx context.Context, eventID, userID uuid.UUID, content string) (entity.Message, error)
}

type ParticipantRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
}

type AssignmentRepository interface {
	Create(ctx context.Context, assignment entity.Assignment) (entity.Assignment, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error)
	DeleteByEvent(ctx context.Context, eventID uuid.UUID) error
	TransactionalDraw(ctx context.Context, eventID uuid.UUID, assignments []entity.Assignment, newStatus definitions.EventStatus) error
}

type ChatRepository interface {
	CreateMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetMessagesByPair(ctx context.Context, eventID, user1ID, user2ID uuid.UUID) ([]entity.Message, error)
}

type EmailService interface {
	SendLoginNotification(ctx context.Context, email, name string) error
	SendOTP(ctx context.Context, email string) (string, error)
	SendDrawNotification(ctx context.Context, email, eventTitle string) error
	SendInvitationEmail(ctx context.Context, email, eventTitle, inviteURL string) error
}

type VerificationRepository interface {
	SaveCode(ctx context.Context, email, code string, expiresAt time.Time) error
	GetValidCode(ctx context.Context, email, code string) (bool, error)
	MarkAsUsed(ctx context.Context, email, code string) error
}
