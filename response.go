package apns2

const (
	APNSErrorPayloadEmpty              = "PayloadEmpty"
	APNSErrorPayloadTooLarge           = "PayloadTooLarge"
	APNSErrorBadTopic                  = "BadTopic"
	APNSErrorTopicDisallowed           = "TopicDisallowed"
	APNSErrorBadMessageId              = "BadMessageId"
	APNSErrorBadExpirationDate         = "BadExpirationDate"
	APNSErrorBadPriority               = "BadPriority"
	APNSErrorMissingDeviceToken        = "MissingDeviceToken"
	APNSErrorBadDeviceToken            = "BadDeviceToken"
	APNSErrorDeviceTokenNotForTopic    = "DeviceTokenNotForTopic"
	APNSErrorUnregistered              = "Unregistered"
	APNSErrorDuplicateHeaders          = "DuplicateHeaders"
	APNSErrorBadCertificateEnvironment = "BadCertificateEnvironment"
	APNSErrorBadCertificate            = "BadCertificate"
	APNSErrorForbidden                 = "Forbidden"
	APNSErrorBadPath                   = "BadPath"
	APNSErrorMethodNotAllowed          = "MethodNotAllowed"
	APNSErrorTooManyRequests           = "TooManyRequests"
	APNSErrorIdleTimeout               = "IdleTimeout"
	APNSErrorShutdown                  = "Shutdown"
	APNSErrorInternalServerError       = "InternalServerError"
	APNSErrorServiceUnavailable        = "ServiceUnavailable"
	APNSErrorMissingTopic              = "MissingTopic"
)

type APNSError struct {
	Reason    string
	Timestamp string
}

func (e *APNSError) Error() string {
	return e.Reason
}

type Response struct {
	StatusCode     int
	NotificationID string
}
