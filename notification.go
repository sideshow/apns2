package apns2

import "time"

const (
	PriorityLow  = 5
	PriorityHigh = 10
)

type Notification struct {
	ApnsId      string
	DeviceToken string
	Topic       string
	Expiration  time.Time
	Priority    int
	Payload     []byte
}
