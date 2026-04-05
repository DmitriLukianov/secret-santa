package response

import (
	"time"

	"secret-santa-backend/internal/entity"
)

type NotificationResponse struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Payload   map[string]string `json:"payload"`
	IsRead    bool              `json:"isRead"`
	CreatedAt time.Time         `json:"createdAt"`
}

func NotificationToResponse(n *entity.Notification) NotificationResponse {
	if n == nil {
		return NotificationResponse{}
	}
	return NotificationResponse{
		ID:        n.ID.String(),
		Type:      n.Type,
		Payload:   n.Payload,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}

func NotificationsToResponse(ns []entity.Notification) []NotificationResponse {
	if ns == nil {
		return nil
	}
	resp := make([]NotificationResponse, len(ns))
	for i := range ns {
		resp[i] = NotificationToResponse(&ns[i])
	}
	return resp
}
