package response

import (
	"time"

	"secret-santa-backend/internal/entity"
)

type FriendshipResponse struct {
	ID          string    `json:"id"`
	RequesterID string    `json:"requesterId"`
	AddresseeID string    `json:"addresseeId"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func FriendshipToResponse(f *entity.Friendship) FriendshipResponse {
	if f == nil {
		return FriendshipResponse{}
	}
	return FriendshipResponse{
		ID:          f.ID.String(),
		RequesterID: f.RequesterID.String(),
		AddresseeID: f.AddresseeID.String(),
		Status:      f.Status,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

func FriendshipsToResponse(fs []entity.Friendship) []FriendshipResponse {
	if fs == nil {
		return nil
	}
	resp := make([]FriendshipResponse, len(fs))
	for i := range fs {
		resp[i] = FriendshipToResponse(&fs[i])
	}
	return resp
}
