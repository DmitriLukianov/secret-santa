package definitions

type EventStatus string

const (
	EventStatusDraft              EventStatus = "draft"
	EventStatusInvitationOpen     EventStatus = "invitation_open"
	EventStatusRegistrationClosed EventStatus = "registration_closed"
	EventStatusDrawingPending     EventStatus = "drawing_pending"
	EventStatusDrawingDone        EventStatus = "drawing_done"
	EventStatusActive             EventStatus = "active"
	EventStatusFinished           EventStatus = "finished"
	EventStatusCancelled          EventStatus = "cancelled"
)

const (
	ParticipantRoleOrganizer   = "organizer"
	ParticipantRoleParticipant = "participant"
)

const (
	WishlistVisibilityPublic    = "public"
	WishlistVisibilityFriends   = "friends"
	WishlistVisibilitySantaOnly = "santa_only"
)

const (
	FriendshipStatusPending  = "pending"
	FriendshipStatusAccepted = "accepted"
	FriendshipStatusDeclined = "declined"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)
