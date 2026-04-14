package definitions

type EventStatus string

const (
	EventStatusRegistration EventStatus = "registration"
	EventStatusGifting      EventStatus = "gifting"
	EventStatusFinished     EventStatus = "finished"
)

const (
	ParticipantRoleOrganizer   = "organizer"
	ParticipantRoleParticipant = "participant"
)

const (
	WishlistVisibilityPublic    = "public"
	WishlistVisibilitySantaOnly = "santa_only"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)
