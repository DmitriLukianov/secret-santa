package request

type SendFriendRequestRequest struct {
	AddresseeID string `json:"addresseeId" validate:"required,uuid"`
}
