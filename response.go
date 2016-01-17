package apns2

import (
	"net/http"
	"strconv"
	"time"
)

// StatusSent is a 200 response.
const StatusSent = http.StatusOK

// The possible Reason error codes returned from APNs.
// From table 6-6 in the Apple Local and Remote Notification Programming Guide.
const (
	ReasonPayloadEmpty              = "PayloadEmpty"
	ReasonPayloadTooLarge           = "PayloadTooLarge"
	ReasonBadTopic                  = "BadTopic"
	ReasonTopicDisallowed           = "TopicDisallowed"
	ReasonBadMessageID              = "BadMessageId"
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

// Response represents a result from the APNs gateway indicating whether a
// notification was accepted or rejected and (if applicable) the metadata
// surrounding the rejection.
type Response struct {

	// The HTTP status code retuened by APNs.
	// A 200 value indicates that the notification was succesfull sent.
	// For a list of other possible status codes, see table 6-4 in the Apple Local
	// and Remote Notification Programming Guide
	StatusCode int

	// The APNs error string indicating the reason for the notification failure (if
	// any). The error code is specified as a string. For a list of possible
	// values, see the Reason constants above.
	// If the notification was accepted, this value will be ""
	Reason string

	// The APNs ApnsID value from the Notification. If you didnt set an ApnsID on the
	// Notification, this will be a new unique UUID whcih has been created by APNs.
	ApnsID string

	// If the value of StatusCode is 410, this is the last time at which APNs
	// confirmed that the device token was no longer valid for the topic.
	Timestamp timestamp
}

// Sent returns whether the notification was succesfull sent.
// This is the same as checking if the StatusCode == 200
func (c *Response) Sent() bool {
	return c.StatusCode == StatusSent
}
