package request

// для жеребьёвки тело не нужно,
// но оставим структуру на будущее

type AssignmentRequest struct {
	EventID string `json:"event_id"`
}
