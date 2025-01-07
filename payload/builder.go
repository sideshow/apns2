// Package payload is a helper package which contains a payload
// builder to make constructing notification payloads easier.
package payload

import "encoding/json"

// InterruptionLevel defines the value for the payload aps interruption-level
type EInterruptionLevel string

const (
	// InterruptionLevelPassive is used to indicate that notification be delivered in a passive manner.
	InterruptionLevelPassive EInterruptionLevel = "passive"

	// InterruptionLevelActive is used to indicate the importance and delivery timing of a notification.
	InterruptionLevelActive EInterruptionLevel = "active"

	// InterruptionLevelTimeSensitive is used to indicate the importance and delivery timing of a notification.
	InterruptionLevelTimeSensitive EInterruptionLevel = "time-sensitive"

	// InterruptionLevelCritical is used to indicate the importance and delivery timing of a notification.
	// This interruption level requires an approved entitlement from Apple.
	// See: https://developer.apple.com/documentation/usernotifications/unnotificationinterruptionlevel/
	InterruptionLevelCritical EInterruptionLevel = "critical"
)

// LiveActivityEvent defines the value for the payload aps event
type ELiveActivityEvent string

const (
	// LiveActivityEventUpdate is used to update an live activity.
	LiveActivityEventUpdate ELiveActivityEvent = "update"

	// LiveActivityEventEnd is used to end an live activity.
	LiveActivityEventEnd ELiveActivityEvent = "end"
)

// Payload represents a notification which holds the content that will be
// marshalled as JSON.
type Payload struct {
	content map[string]interface{}
}

type aps struct {
	Alert             interface{}            `json:"alert,omitempty"`
	Badge             interface{}            `json:"badge,omitempty"`
	Category          string                 `json:"category,omitempty"`
	ContentAvailable  int                    `json:"content-available,omitempty"`
	InterruptionLevel EInterruptionLevel     `json:"interruption-level,omitempty"`
	MutableContent    int                    `json:"mutable-content,omitempty"`
	RelevanceScore    interface{}            `json:"relevance-score,omitempty"`
	Sound             interface{}            `json:"sound,omitempty"`
	ThreadID          string                 `json:"thread-id,omitempty"`
	URLArgs           []string               `json:"url-args,omitempty"`
	ContentState      map[string]interface{} `json:"content-state,omitempty"`
	DismissalDate     int64                  `json:"dismissal-date,omitempty"`
	StaleDate         int64                  `json:"stale-date,omitempty"`
	Event             ELiveActivityEvent     `json:"event,omitempty"`
	Timestamp         int64                  `json:"timestamp,omitempty"`
	AttributesType    string                 `json:"attributes-type,omitempty"`
	Attributes        map[string]interface{} `json:"attributes,omitempty"`
}

type alert struct {
	Action          string   `json:"action,omitempty"`
	ActionLocKey    string   `json:"action-loc-key,omitempty"`
	Body            string   `json:"body,omitempty"`
	LaunchImage     string   `json:"launch-image,omitempty"`
	LocArgs         []string `json:"loc-args,omitempty"`
	LocKey          string   `json:"loc-key,omitempty"`
	Title           string   `json:"title,omitempty"`
	Subtitle        string   `json:"subtitle,omitempty"`
	TitleLocArgs    []string `json:"title-loc-args,omitempty"`
	TitleLocKey     string   `json:"title-loc-key,omitempty"`
	SubtitleLocArgs []string `json:"subtitle-loc-args,omitempty"`
	SubtitleLocKey  string   `json:"subtitle-loc-key,omitempty"`
	SummaryArg      string   `json:"summary-arg,omitempty"`
	SummaryArgCount int      `json:"summary-arg-count,omitempty"`
}

type sound struct {
	Critical int     `json:"critical,omitempty"`
	Name     string  `json:"name,omitempty"`
	Volume   float32 `json:"volume,omitempty"`
}

// NewPayload returns a new Payload struct
func NewPayload() *Payload {
	return &Payload{
		map[string]interface{}{
			"aps": &aps{},
		},
	}
}

// Alert sets the aps alert on the payload.
// This will display a notification alert message to the user.
//
//	{"aps":{"alert":alert}}`
func (p *Payload) Alert(alert interface{}) *Payload {
	p.aps().Alert = alert
	return p
}

// SetContentState sets the aps content-state on the payload.
// This will update content-state of live activity widget.
//
//	{"aps":{"content-state": {} }}`
func (p *Payload) SetContentState(contentState map[string]interface{}) *Payload {
	p.aps().ContentState = contentState
	return p
}

// SetDismissalDate sets the aps dismissal-date on the payload.
// This will remove the live activity from the user's UI at the given timestamp.
//
//	{"aps":{"dismissal-date": DismissalDate }}`
func (p *Payload) SetDismissalDate(dismissalDate int64) *Payload {
	p.aps().DismissalDate = dismissalDate
	return p
}

// SetStaleDate sets the aps stale-date on the payload.
// This will mark this live activity update as outdated at the given timestamp.
//
//	{"aps":{"stale-date": StaleDate }}`
func (p *Payload) SetStaleDate(staleDate int64) *Payload {
	p.aps().StaleDate = staleDate
	return p
}

// SetEvent sets the aps event type on the payload.
// This can either be `LiveActivityEventUpdate` or `LiveActivityEventEnd`
//
//	{"aps":{"event": Event }}`
func (p *Payload) SetEvent(event ELiveActivityEvent) *Payload {
	p.aps().Event = event
	return p
}

// SetTimestamp sets the aps timestamp on the payload.
// This will let live activity know when to update the stuff.
//
//	{"aps":{"timestamp": Timestamp }}`
func (p *Payload) SetTimestamp(timestamp int64) *Payload {
	p.aps().Timestamp = timestamp
	return p
}

// SetAttributesType sets the aps attributes-type field on the payload.
// This is used for push-to-start live activities
//
//	{"aps":{"attributes-type": attributesType }}`
func (p *Payload) SetAttributesType(attributesType string) *Payload {
	p.aps().AttributesType = attributesType
	return p
}

// SetAttributes sets the aps attributes field on the payload.
// This is used for push-to-start live activities
//
//	{"aps":{"attributes": attributes }}`
func (p *Payload) SetAttributes(attributes map[string]interface{}) *Payload {
	p.aps().Attributes = attributes
	return p
}

// Badge sets the aps badge on the payload.
// This will display a numeric badge on the app icon.
//
//	{"aps":{"badge":b}}
func (p *Payload) Badge(b int) *Payload {
	p.aps().Badge = b
	return p
}

// ZeroBadge sets the aps badge on the payload to 0.
// This will clear the badge on the app icon.
//
//	{"aps":{"badge":0}}
func (p *Payload) ZeroBadge() *Payload {
	p.aps().Badge = 0
	return p
}

// UnsetBadge removes the badge attribute from the payload.
// This will leave the badge on the app icon unchanged.
// If you wish to clear the app icon badge, use ZeroBadge() instead.
//
//	{"aps":{}}
func (p *Payload) UnsetBadge() *Payload {
	p.aps().Badge = nil
	return p
}

// Sound sets the aps sound on the payload.
// This will play a sound from the app bundle, or the default sound otherwise.
//
//	{"aps":{"sound":sound}}
func (p *Payload) Sound(sound interface{}) *Payload {
	p.aps().Sound = sound
	return p
}

// ContentAvailable sets the aps content-available on the payload to 1.
// This will indicate to the app that there is new content available to download
// and launch the app in the background.
//
//	{"aps":{"content-available":1}}
func (p *Payload) ContentAvailable() *Payload {
	p.aps().ContentAvailable = 1
	return p
}

// MutableContent sets the aps mutable-content on the payload to 1.
// This will indicate to the to the system to call your Notification Service
// extension to mutate or replace the notification's content.
//
//	{"aps":{"mutable-content":1}}
func (p *Payload) MutableContent() *Payload {
	p.aps().MutableContent = 1
	return p
}

// Custom payload

// Custom sets a custom key and value on the payload.
// This will add custom key/value data to the notification payload at root level.
//
//	{"aps":{}, key:value}
func (p *Payload) Custom(key string, val interface{}) *Payload {
	p.content[key] = val
	return p
}

// UnsetCustom unsets a custom key and value on the payload.
// This will delete custom key/value data from the notification payload at root level.
//
//	{"aps":{}}
func (p *Payload) UnsetCustom(key string) *Payload {
	delete(p.content, key)
	return p
}

// Alert dictionary

// AlertTitle sets the aps alert title on the payload.
// This will display a short string describing the purpose of the notification.
// Apple Watch & Safari display this string as part of the notification interface.
//
//	{"aps":{"alert":{"title":title}}}
func (p *Payload) AlertTitle(title string) *Payload {
	p.aps().alert().Title = title
	return p
}

// AlertTitleLocKey sets the aps alert title localization key on the payload.
// This is the key to a title string in the Localizable.strings file for the
// current localization. See Localized Formatted Strings in Apple documentation
// for more information.
//
//	{"aps":{"alert":{"title-loc-key":key}}}
func (p *Payload) AlertTitleLocKey(key string) *Payload {
	p.aps().alert().TitleLocKey = key
	return p
}

// AlertTitleLocArgs sets the aps alert title localization args on the payload.
// These are the variable string values to appear in place of the format
// specifiers in title-loc-key. See Localized Formatted Strings in Apple
// documentation for more information.
//
//	{"aps":{"alert":{"title-loc-args":args}}}
func (p *Payload) AlertTitleLocArgs(args []string) *Payload {
	p.aps().alert().TitleLocArgs = args
	return p
}

// AlertSubtitle sets the aps alert subtitle on the payload.
// This will display a short string describing the purpose of the notification.
// Apple Watch & Safari display this string as part of the notification interface.
//
//	{"aps":{"alert":{"subtitle":"subtitle"}}}
func (p *Payload) AlertSubtitle(subtitle string) *Payload {
	p.aps().alert().Subtitle = subtitle
	return p
}

// AlertSubtitleLocKey sets the aps alert subtitle localization key on the payload.
// This is the key to a subtitle string in the Localizable.strings file for the
// current localization. See Localized Formatted Strings in Apple documentation
// for more information.
//
//	{"aps":{"alert":{"subtitle-loc-key":key}}}
func (p *Payload) AlertSubtitleLocKey(key string) *Payload {
	p.aps().alert().SubtitleLocKey = key
	return p
}

// AlertSubtitleLocArgs sets the aps alert subtitle localization args on the payload.
// These are the variable string values to appear in place of the format
// specifiers in subtitle-loc-key. See Localized Formatted Strings in Apple
// documentation for more information.
//
//	{"aps":{"alert":{"title-loc-args":args}}}
func (p *Payload) AlertSubtitleLocArgs(args []string) *Payload {
	p.aps().alert().SubtitleLocArgs = args
	return p
}

// AlertBody sets the aps alert body on the payload.
// This is the text of the alert message.
//
//	{"aps":{"alert":{"body":body}}}
func (p *Payload) AlertBody(body string) *Payload {
	p.aps().alert().Body = body
	return p
}

// AlertLaunchImage sets the aps launch image on the payload.
// This is the filename of an image file in the app bundle. The image is used
// as the launch image when users tap the action button or move the action
// slider.
//
//	{"aps":{"alert":{"launch-image":image}}}
func (p *Payload) AlertLaunchImage(image string) *Payload {
	p.aps().alert().LaunchImage = image
	return p
}

// AlertLocArgs sets the aps alert localization args on the payload.
// These are the variable string values to appear in place of the format
// specifiers in loc-key. See Localized Formatted Strings in Apple
// documentation for more information.
//
//	{"aps":{"alert":{"loc-args":args}}}
func (p *Payload) AlertLocArgs(args []string) *Payload {
	p.aps().alert().LocArgs = args
	return p
}

// AlertLocKey sets the aps alert localization key on the payload.
// This is the key to an alert-message string in the Localizable.strings file
// for the current localization. See Localized Formatted Strings in Apple
// documentation for more information.
//
//	{"aps":{"alert":{"loc-key":key}}}
func (p *Payload) AlertLocKey(key string) *Payload {
	p.aps().alert().LocKey = key
	return p
}

// AlertAction sets the aps alert action on the payload.
// This is the label of the action button, if the user sets the notifications
// to appear as alerts. This label should be succinct, such as “Details” or
// “Read more”. If omitted, the default value is “Show”.
//
//	{"aps":{"alert":{"action":action}}}
func (p *Payload) AlertAction(action string) *Payload {
	p.aps().alert().Action = action
	return p
}

// AlertActionLocKey sets the aps alert action localization key on the payload.
// This is the the string used as a key to get a localized string in the current
// localization to use for the notfication right button’s title instead of
// “View”. See Localized Formatted Strings in Apple documentation for more
// information.
//
//	{"aps":{"alert":{"action-loc-key":key}}}
func (p *Payload) AlertActionLocKey(key string) *Payload {
	p.aps().alert().ActionLocKey = key
	return p
}

// AlertSummaryArg sets the aps alert summary arg key on the payload.
// This is the string that is used as a key to fill in an argument
// at the bottom of a notification to provide more context, such as
// a name associated with the sender of the notification.
//
//	{"aps":{"alert":{"summary-arg":key}}}
func (p *Payload) AlertSummaryArg(key string) *Payload {
	p.aps().alert().SummaryArg = key
	return p
}

// AlertSummaryArgCount sets the aps alert summary arg count key on the payload.
// This integer sets a custom "weight" on the notification, effectively
// allowing a notification to be viewed internally as two. For example if
// a notification encompasses 3 messages, you can set it to 3.
//
//	{"aps":{"alert":{"summary-arg-count":key}}}
func (p *Payload) AlertSummaryArgCount(key int) *Payload {
	p.aps().alert().SummaryArgCount = key
	return p
}

// General

// Category sets the aps category on the payload.
// This is a string value that represents the identifier property of the
// UIMutableUserNotificationCategory object you created to define custom actions.
//
//	{"aps":{"category":category}}
func (p *Payload) Category(category string) *Payload {
	p.aps().Category = category
	return p
}

// Mdm sets the mdm on the payload.
// This is for Apple Mobile Device Management (mdm) payloads.
//
//	{"aps":{}:"mdm":mdm}
func (p *Payload) Mdm(mdm string) *Payload {
	p.content["mdm"] = mdm
	return p
}

// ThreadID sets the aps thread id on the payload.
// This is for the purpose of updating the contents of a View Controller in a
// Notification Content app extension when a new notification arrives. If a
// new notification arrives whose thread-id value matches the thread-id of the
// notification already being displayed, the didReceiveNotification method
// is called.
//
//	{"aps":{"thread-id":id}}
func (p *Payload) ThreadID(threadID string) *Payload {
	p.aps().ThreadID = threadID
	return p
}

// URLArgs sets the aps category on the payload.
// This specifies an array of values that are paired with the placeholders
// inside the urlFormatString value of your website.json file.
// See Apple Notification Programming Guide for Websites.
//
//	{"aps":{"url-args":urlArgs}}
func (p *Payload) URLArgs(urlArgs []string) *Payload {
	p.aps().URLArgs = urlArgs
	return p
}

// SoundName sets the name value on the aps sound dictionary.
// This function makes the notification a critical alert, which should be pre-approved by Apple.
// See: https://developer.apple.com/contact/request/notifications-critical-alerts-entitlement/
//
// {"aps":{"sound":{"critical":1,"name":name,"volume":1.0}}}
func (p *Payload) SoundName(name string) *Payload {
	p.aps().sound().Name = name
	return p
}

// SoundVolume sets the volume value on the aps sound dictionary.
// This function makes the notification a critical alert, which should be pre-approved by Apple.
// See: https://developer.apple.com/contact/request/notifications-critical-alerts-entitlement/
//
// {"aps":{"sound":{"critical":1,"name":"default","volume":volume}}}
func (p *Payload) SoundVolume(volume float32) *Payload {
	p.aps().sound().Volume = volume
	return p
}

// InterruptionLevel defines the value for the payload aps interruption-level
// This is to indicate the importance and delivery timing of a notification.
// (Using InterruptionLevelCritical requires an approved entitlement from Apple.)
// See: https://developer.apple.com/documentation/usernotifications/unnotificationinterruptionlevel/
//
// {"aps":{"interruption-level":passive}}
func (p *Payload) InterruptionLevel(interruptionLevel EInterruptionLevel) *Payload {
	p.aps().InterruptionLevel = interruptionLevel
	return p
}

// The relevance score, a number between 0 and 1,
// that the system uses to sort the notifications from your app.
// The highest score gets featured in the notification summary.
// See https://developer.apple.com/documentation/usernotifications/unnotificationcontent/3821031-relevancescore.
//
//	{"aps":{"relevance-score":0.1}}
func (p *Payload) RelevanceScore(b float32) *Payload {
	p.aps().RelevanceScore = b
	return p
}

// Unsets the relevance score
// that the system uses to sort the notifications from your app.
// The highest score gets featured in the notification summary.
// See https://developer.apple.com/documentation/usernotifications/unnotificationcontent/3821031-relevancescore.
//
//	{"aps":{"relevance-score":0.1}}
func (p *Payload) UnsetRelevanceScore() *Payload {
	p.aps().RelevanceScore = nil
	return p
}

// MarshalJSON returns the JSON encoded version of the Payload
func (p *Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.content)
}

func (p *Payload) aps() *aps {
	return p.content["aps"].(*aps)
}

func (a *aps) alert() *alert {
	if _, ok := a.Alert.(*alert); !ok {
		a.Alert = &alert{}
	}
	return a.Alert.(*alert)
}

func (a *aps) sound() *sound {
	if _, ok := a.Sound.(*sound); !ok {
		a.Sound = &sound{Critical: 1, Name: "default", Volume: 1.0}
	}
	return a.Sound.(*sound)
}
