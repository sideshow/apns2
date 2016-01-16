package apns2_test

import (
	"encoding/json"
	"net/http"
	"testing"

	apns "github.com/sideshow/apns2"
)

// Unit Tests

func TestResponseSent(t *testing.T) {
	if apns.StatusSent != http.StatusOK {
		t.Error("StatusSent should be", http.StatusOK)
	}
	if (&apns.Response{StatusCode: 200}).Sent() == false {
		t.Error("Sent() should be true")
	}
	if (&apns.Response{StatusCode: 400}).Sent() {
		t.Error("Sent() should be false")
	}
}

func TestStringTimestampParse(t *testing.T) {
	response := &apns.Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\":\"1421147681\"}"
	json.Unmarshal([]byte(payload), &response)
	if response.Timestamp.Unix() != 1421147681 {
		t.Error("Timestamp should be", 1421147681)
	}
}

func TestIntTimestampParse(t *testing.T) {
	response := &apns.Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\":1421147681}"
	json.Unmarshal([]byte(payload), &response)
	if response.Timestamp.Unix() != 1421147681 {
		t.Error("Timestamp should be", 1421147681)
	}
}

func TestInvalidTimestampParse(t *testing.T) {
	response := &apns.Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\": \"2016-01-16 17:44:04 +1300\"}"
	err := json.Unmarshal([]byte(payload), &response)
	if err == nil {
		t.Error("Bad Timestamp should generate parsing error")
	}
}
