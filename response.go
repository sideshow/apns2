package apns2

import (
	"net/http"
	"strconv"
	"time"
)

const StatusSent = http.StatusOK
const (
	ReasonPayloadEmpty              = "PayloadEmpty"
	ReasonPayloadTooLarge           = "PayloadTooLarge"
	ReasonBadTopic                  = "BadTopic"
	ReasonTopicDisallowed           = "TopicDisallowed"
	ReasonBadMessageId              = "BadMessageId"
	ReasonBadExpirationDate         = "BadExpirationDate"
	ReasonBadPriority               = "BadPriority"
	ReasonMissingDeviceToken        = "MissingDeviceToken"
	ReasonBadDeviceToken            = "BadDeviceToken"
	ReasonDeviceTokenNotForTopic    = "DeviceTokenNotForTopic"
	ReasonUnregistered              = "Unregistered"
	ReasonDuplicateHeaders          = "DuplicateHeaders"
	ReasonBadCertificateEnvironment = "BadCertificateEnvironment"
	ReasonBadCertificate            = "BadCertificate"
	ReasonForbidden                 = "Forbidden"
	ReasonBadPath                   = "BadPath"
	ReasonMethodNotAllowed          = "MethodNotAllowed"
	ReasonTooManyRequests           = "TooManyRequests"
	ReasonIdleTimeout               = "IdleTimeout"
	ReasonShutdown                  = "Shutdown"
	ReasonInternalServerError       = "InternalServerError"
	ReasonServiceUnavailable        = "ServiceUnavailable"
	ReasonMissingTopic              = "MissingTopic"
)

type timestamp struct {
	time.Time
}

func (t *timestamp) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

type Response struct {
	StatusCode int
	Reason     string
	ApnsId     string
	Timestamp  timestamp
}

func (c *Response) Sent() bool {
	return c.StatusCode == StatusSent
}
