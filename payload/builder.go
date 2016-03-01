// Package payload is a helper package which contains a payload
// builder to make constructing notification payloads easier.
package payload

import "encoding/json"

type payload struct {
	content map[string]interface{}
}

type aps struct {
	Alert            interface{} `json:"alert,omitempty"`
	Badge            interface{} `json:"badge,omitempty"`
	Category         string      `json:"category,omitempty"`
	ContentAvailable int         `json:"content-available,omitempty"`
	URLArgs          []string    `json:"url-args,omitempty"`
	Sound            string      `json:"sound,omitempty"`
}

type alert struct {
	Action       string   `json:"action,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	Body         string   `json:"body,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	Title        string   `json:"title,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
}

// NewPayload represents a notification payload.
func NewPayload() *payload {
	return &payload{
		map[string]interface{}{
			"aps": &aps{},
		},
	}
}

// Sets the aps alert on the payload.
// This will display a notification alert message to the user.
// {"aps":{"alert":alert}}
func (p *payload) Alert(alert interface{}) *payload {
	p.aps().Alert = alert
	return p
}

// Sets the aps badge on the payload.
// This will display a numeric badge on the app icon.
// {"aps":{"badge":b}}
func (p *payload) Badge(b int) *payload {
	p.aps().Badge = b
	return p
}

// Sets the aps badge on the payload to 0.
// This will clear the badge on the app icon.
// {"aps":{"badge":0}}
func (p *payload) ZeroBadge() *payload {
	p.aps().Badge = 0
	return p
}

// Removes the badge attribute from the payload.
// This will leave the badge on the app icon unchanged.
// If you wish to clear the app icon badge, use ZeroBadge() instead.
// {"aps":{}}
func (p *payload) UnsetBadge() *payload {
	p.aps().Badge = nil
	return p
}

// Sets the aps sound on the payload.
// This will play a sound from the app bundle, or the default sound otherwise.
// {"aps":{"sound":sound}}
func (p *payload) Sound(sound string) *payload {
	p.aps().Sound = sound
	return p
}

// Sets the aps content-available on the payload to 1.
// This will indicate to the app that there is new content available to download
// and launch the app in the background.
// {"aps":{"content-available":1}}
func (p *payload) ContentAvailable() *payload {
	p.aps().ContentAvailable = 1
	return p
}

// Custom payload

// Sets a custom key and value on the payload.
// This will add custom key/value data to the notification payload at root level.
// {"aps":{}, key:value}
func (p *payload) Custom(key string, val interface{}) *payload {
	p.content[key] = val
	return p
}

// Alert dictionary

// Sets the aps alert title on the payload.
// This will display a short string describing the purpose of the notification.
// Apple Watch & Safari display this string as part of the notification interface.
// {"aps":{"alert":"title"}}
func (p *payload) AlertTitle(title string) *payload {
	p.aps().alert().Title = title
	return p
}

// Sets the aps alert title localization key on the payload.
// This is the key to a title string in the Localizable.strings file for the
// current localization. See Localized Formatted Strings in Apple documentation
// for more information.
// {"aps":{"alert":{"title-loc-key":key}}}
func (p *payload) AlertTitleLocKey(key string) *payload {
	p.aps().alert().TitleLocKey = key
	return p
}

// Sets the aps alert title localization args on the payload.
// These are the variable string values to appear in place of the format
// specifiers in title-loc-key. See Localized Formatted Strings in Apple
// documentation for more information.
// {"aps":{"alert":{"title-loc-args":args}}}
func (p *payload) AlertTitleLocArgs(args []string) *payload {
	p.aps().alert().TitleLocArgs = args
	return p
}

// Sets the aps alert body on the payload.
// This is the text of the alert message.
// {"aps":{"alert":{"body":body}}}
func (p *payload) AlertBody(body string) *payload {
	p.aps().alert().Body = body
	return p
}

// Sets the aps launch image on the payload.
// This is the filename of an image file in the app bundle. The image is used
// as the launch image when users tap the action button or move the action
// slider.
// {"aps":{"alert":{"launch-image":image}}}
func (p *payload) AlertLaunchImage(image string) *payload {
	p.aps().alert().LaunchImage = image
	return p
}

// Sets the aps alert localization key on the payload.
// This is the key to an alert-message string in the Localizable.strings file
// for the current localization. See Localized Formatted Strings in Apple
// documentation for more information.
// {"aps":{"alert":{"loc-key":key}}}
func (p *payload) AlertLocKey(key string) *payload {
	p.aps().alert().LocKey = key
	return p
}

// Sets the aps alert action on the payload.
// This is the label of the action button, if the user sets the notifications
// to appear as alerts. This label should be succinct, such as “Details” or
// “Read more”. If omitted, the default value is “Show”.
// {"aps":{"alert":{"action":action}}}
func (p *payload) AlertAction(action string) *payload {
	p.aps().alert().Action = action
	return p
}

// Sets the aps alert action localization key on the payload.
// This is the the string used as a key to get a localized string in the current
// localization to use for the notfication right button’s title instead of
// “View”. See Localized Formatted Strings in Apple documentation for more
// information.
// {"aps":{"alert":{"action-loc-key":key}}}
func (p *payload) AlertActionLocKey(key string) *payload {
	p.aps().alert().ActionLocKey = key
	return p
}

// General

// Sets the aps category on the payload.
// This is a string value that represents the identifier property of the
// UIMutableUserNotificationCategory object you created to define custom actions.
// {"aps":{"alert":{"category":category}}}
func (p *payload) Category(category string) *payload {
	p.aps().Category = category
	return p
}

// Sets the mdm on the payload.
// This is for Apple Mobile Device Management (mdm) payloads.
// {"aps":{}:"mdm":mdm}
func (p *payload) Mdm(mdm string) *payload {
	p.content["mdm"] = mdm
	return p
}

// Sets the aps category on the payload.
// This specifies an array of values that are paired with the placeholders
// inside the urlFormatString value of your website.json file.
// See Apple Notification Programming Guide for Websites.
// {"aps":{"url-args":urlArgs}}
func (p *payload) URLArgs(urlArgs []string) *payload {
	p.aps().URLArgs = urlArgs
	return p
}

func (p *payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.content)
}

func (p *payload) aps() *aps {
	return p.content["aps"].(*aps)
}

func (a *aps) alert() *alert {
	if _, ok := a.Alert.(*alert); !ok {
		a.Alert = &alert{}
	}
	return a.Alert.(*alert)
}
