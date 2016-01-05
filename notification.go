package apns2

import "time"

const (
	PriorityLow  = "5"
	PriorityHigh = "10"
)

type Notification struct {
	Id          string
	DeviceToken string
	Topic       string
	Expiry      time.Time
	Priority    int
	Payload     []byte
}
