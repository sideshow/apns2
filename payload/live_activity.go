package payload

type ELiveActivityEvent string

const (
	LiveActivityEventUpdate ELiveActivityEvent = "update"
	LiveActivityEventEnd    ELiveActivityEvent = "end"
)
