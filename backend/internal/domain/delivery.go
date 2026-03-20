package domain

import "time"

type Delivery struct {
	ID             string
	EventID        string
	SenderID       string
	ReceiverID     string
	DeliveryType   string
	TrackingNumber string
	Status         string
	CreatedAt      time.Time
}
