package definitions

// ====================== EVENT STATUS ======================
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

// ====================== PARTICIPANT ROLE ======================
const (
	ParticipantRoleOrganizer   = "organizer"
	ParticipantRoleParticipant = "participant"
)

// ====================== WISHLIST VISIBILITY ======================
const (
	WishlistVisibilityPublic    = "public"
	WishlistVisibilityFriends   = "friends"
	WishlistVisibilitySantaOnly = "santa_only"
)
