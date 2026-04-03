package response

import "secret-santa-backend/internal/entity"

func EventToResponse(e *entity.Event) EventResponse {
	if e == nil {
		return EventResponse{}
	}
	return EventResponse{
		ID:              e.ID.String(),
		Title:           e.Title,
		Description:     e.Description,
		Rules:           e.Rules,
		Recommendations: e.Recommendations,
		OrganizerID:     e.OrganizerID.String(),
		StartDate:       e.StartDate,
		DrawDate:        e.DrawDate,
		EndDate:         e.EndDate,
		Status:          string(e.Status),
		MaxParticipants: e.MaxParticipants,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func EventsToResponse(events []entity.Event) []EventResponse {
	if events == nil {
		return nil
	}

	resp := make([]EventResponse, len(events))
	for i := range events {
		resp[i] = EventToResponse(&events[i])
	}
	return resp
}

func UserToResponse(u *entity.User) UserResponse {
	if u == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
	}
}

func UsersToResponse(users []*entity.User) []UserResponse {
	if users == nil {
		return nil
	}

	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = UserToResponse(u)
	}
	return resp
}

func ParticipantToResponse(p *entity.Participant) ParticipantResponse {
	if p == nil {
		return ParticipantResponse{}
	}
	return ParticipantResponse{
		ID:         p.ID.String(),
		EventID:    p.EventID.String(),
		UserID:     p.UserID.String(),
		Role:       p.Role,
		GiftSent:   p.GiftSent,
		GiftSentAt: p.GiftSentAt,
		CreatedAt:  p.CreatedAt,
	}
}

func ParticipantsToResponse(participants []*entity.Participant) []ParticipantResponse {
	if participants == nil {
		return nil
	}

	resp := make([]ParticipantResponse, len(participants))
	for i, p := range participants {
		resp[i] = ParticipantToResponse(p)
	}
	return resp
}

func WishlistToResponse(w *entity.Wishlist) WishlistResponse {
	if w == nil {
		return WishlistResponse{}
	}
	return WishlistResponse{
		ID:            w.ID.String(),
		ParticipantID: w.ParticipantID.String(),
		Visibility:    w.Visibility,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}

func WishlistItemToResponse(item *entity.WishlistItem) WishlistItemResponse {
	if item == nil {
		return WishlistItemResponse{}
	}
	return WishlistItemResponse{
		ID:        item.ID.String(),
		Title:     item.Title,
		Link:      item.Link,
		ImageURL:  item.ImageURL,
		Comment:   item.Comment,
		CreatedAt: item.CreatedAt,
	}
}

func WishlistItemsToResponse(items []*entity.WishlistItem) []WishlistItemResponse {
	if items == nil {
		return nil
	}

	resp := make([]WishlistItemResponse, len(items))
	for i, item := range items {
		resp[i] = WishlistItemToResponse(item)
	}
	return resp
}

func AssignmentToResponse(a *entity.Assignment) AssignmentResponse {
	if a == nil {
		return AssignmentResponse{}
	}
	return AssignmentResponse{
		ID:         a.ID.String(),
		EventID:    a.EventID.String(),
		GiverID:    a.GiverID.String(),
		ReceiverID: a.ReceiverID.String(),
		CreatedAt:  a.CreatedAt,
	}
}

func AssignmentsToResponse(assignments []*entity.Assignment) []AssignmentResponse {
	if assignments == nil {
		return nil
	}

	resp := make([]AssignmentResponse, len(assignments))
	for i, a := range assignments {
		resp[i] = AssignmentToResponse(a)
	}
	return resp
}
