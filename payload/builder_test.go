package payload_test

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/deltapath-eric/apns2/payload"
	"github.com/stretchr/testify/assert"
)

func TestEmptyPayload(t *testing.T) {
	payload := NewPayload()
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{}}`, string(b))
}

func TestAlert(t *testing.T) {
	payload := NewPayload().Alert("hello")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":"hello"}}`, string(b))
}

func TestBadge(t *testing.T) {
	payload := NewPayload().Badge(1)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"badge":1}}`, string(b))
}

func TestZeroBadge(t *testing.T) {
	payload := NewPayload().ZeroBadge()
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"badge":0}}`, string(b))
}

func TestUnsetBadge(t *testing.T) {
	payload := NewPayload().Badge(1).UnsetBadge()
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{}}`, string(b))
}

func TestSound(t *testing.T) {
	payload := NewPayload().Sound("Default.caf")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"sound":"Default.caf"}}`, string(b))
}

func TestSoundDictionary(t *testing.T) {
	payload := NewPayload().Sound(map[string]any{
		"critical": 1,
		"name":     "default",
		"volume":   0.8,
	})
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"sound":{"critical":1,"name":"default","volume":0.8}}}`, string(b))
}

func TestContentAvailable(t *testing.T) {
	payload := NewPayload().ContentAvailable()
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"content-available":1}}`, string(b))
}

func TestMutableContent(t *testing.T) {
	payload := NewPayload().MutableContent()
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"mutable-content":1}}`, string(b))
}

func TestCustom(t *testing.T) {
	payload := NewPayload().Custom("key", "val")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{},"key":"val"}`, string(b))
}

func TestCustomMap(t *testing.T) {
	payload := NewPayload().Custom("key", map[string]any{
		"map": 1,
	})
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{},"key":{"map":1}}`, string(b))
}

func TestAlertTitle(t *testing.T) {
	payload := NewPayload().AlertTitle("hello")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"title":"hello"}}}`, string(b))
}

func TestAlertTitleLocKey(t *testing.T) {
	payload := NewPayload().AlertTitleLocKey("GAME_PLAY_REQUEST_FORMAT")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"title-loc-key":"GAME_PLAY_REQUEST_FORMAT"}}}`, string(b))
}

func TestAlertLocArgs(t *testing.T) {
	payload := NewPayload().AlertLocArgs([]string{"Jenna", "Frank"})
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"loc-args":["Jenna","Frank"]}}}`, string(b))
}

func TestAlertTitleLocArgs(t *testing.T) {
	payload := NewPayload().AlertTitleLocArgs([]string{"Jenna", "Frank"})
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"title-loc-args":["Jenna","Frank"]}}}`, string(b))
}

func TestAlertSubtitle(t *testing.T) {
	payload := NewPayload().AlertSubtitle("hello")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"subtitle":"hello"}}}`, string(b))
}

func TestAlertBody(t *testing.T) {
	payload := NewPayload().AlertBody("body")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"body":"body"}}}`, string(b))
}

func TestAlertLaunchImage(t *testing.T) {
	payload := NewPayload().AlertLaunchImage("Default.png")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"launch-image":"Default.png"}}}`, string(b))
}

func TestAlertLocKey(t *testing.T) {
	payload := NewPayload().AlertLocKey("LOC")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"loc-key":"LOC"}}}`, string(b))
}

func TestAlertAction(t *testing.T) {
	payload := NewPayload().AlertAction("action")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"action":"action"}}}`, string(b))
}

func TestAlertActionLocKey(t *testing.T) {
	payload := NewPayload().AlertActionLocKey("PLAY")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"action-loc-key":"PLAY"}}}`, string(b))
}

func TestCategory(t *testing.T) {
	payload := NewPayload().Category("NEW_MESSAGE_CATEGORY")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"category":"NEW_MESSAGE_CATEGORY"}}`, string(b))
}

func TestContentState(t *testing.T) {
	payload := NewPayload().SetContentState(map[string]any{"my_int": 13, "my_string": "foo"})
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"content-state":{"my_int":13,"my_string":"foo"}}}`, string(b))
}

func TestDismissalDate(t *testing.T) {
	timestamp := time.Date(2023, 1, 27, 12, 14, 00, 00, time.UTC).Unix()
	payload := NewPayload().SetDismissalDate(timestamp)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"dismissal-date":1674821640}}`, string(b))
}

func TestStaleDate(t *testing.T) {
	timestamp := time.Date(2023, 1, 27, 12, 14, 00, 00, time.UTC).Unix()
	payload := NewPayload().SetStaleDate(timestamp)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"stale-date":1674821640}}`, string(b))
}

func TestEventEnd(t *testing.T) {
	payload := NewPayload().SetEvent(LiveActivityEventEnd)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"event":"end"}}`, string(b))
}

func TestEventUpdate(t *testing.T) {
	payload := NewPayload().SetEvent(LiveActivityEventUpdate)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"event":"update"}}`, string(b))
}

func TestTimestamp(t *testing.T) {
	timestamp := time.Date(2023, 1, 27, 12, 14, 00, 00, time.UTC).Unix()
	payload := NewPayload().SetTimestamp(timestamp)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"timestamp":1674821640}}`, string(b))
}

func TestMdm(t *testing.T) {
	payload := NewPayload().Mdm("996ac527-9993-4a0a-8528-60b2b3c2f52b")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{},"mdm":"996ac527-9993-4a0a-8528-60b2b3c2f52b"}`, string(b))
}

func TestThreadID(t *testing.T) {
	payload := NewPayload().ThreadID("THREAD_ID")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"thread-id":"THREAD_ID"}}`, string(b))
}

func TestURLArgs(t *testing.T) {
	payload := NewPayload().URLArgs([]string{"a", "b"})
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"url-args":["a","b"]}}`, string(b))
}

func TestSoundName(t *testing.T) {
	payload := NewPayload().SoundName("test")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"sound":{"critical":1,"name":"test","volume":1}}}`, string(b))
}

func TestSoundVolume(t *testing.T) {
	payload := NewPayload().SoundVolume(0.5)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"sound":{"critical":1,"name":"default","volume":0.5}}}`, string(b))
}

func TestAlertSummaryArg(t *testing.T) {
	payload := NewPayload().AlertSummaryArg("Robert")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"summary-arg":"Robert"}}}`, string(b))
}

func TestAlertSummaryArgCount(t *testing.T) {
	payload := NewPayload().AlertSummaryArgCount(3)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":{"summary-arg-count":3}}}`, string(b))
}

func TestInterruptionLevelPassive(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelPassive)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"interruption-level":"passive"}}`, string(b))
}

func TestInterruptionLevelActive(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelActive)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"interruption-level":"active"}}`, string(b))
}

func TestInterruptionLevelTimeSensitive(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelTimeSensitive)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"interruption-level":"time-sensitive"}}`, string(b))
}

func TestInterruptionLevelCritical(t *testing.T) {
	payload := NewPayload().InterruptionLevel(InterruptionLevelCritical)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"interruption-level":"critical"}}`, string(b))
}

func TestRelevanceScore(t *testing.T) {
	payload := NewPayload().RelevanceScore(0.1)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"relevance-score":0.1}}`, string(b))
}

func TestRelevanceScoreZero(t *testing.T) {
	payload := NewPayload().RelevanceScore(0)
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"relevance-score":0}}`, string(b))
}

func TestUnsetRelevanceScore(t *testing.T) {
	payload := NewPayload().RelevanceScore(0.1).UnsetRelevanceScore()
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{}}`, string(b))
}

func TestCombined(t *testing.T) {
	payload := NewPayload().Alert("hello").Badge(1).Sound("Default.caf").InterruptionLevel(InterruptionLevelActive).RelevanceScore(0.1).Custom("key", "val")
	b, _ := json.Marshal(payload)
	assert.Equal(t, `{"aps":{"alert":"hello","badge":1,"interruption-level":"active","relevance-score":0.1,"sound":"Default.caf"},"key":"val"}`, string(b))
}
