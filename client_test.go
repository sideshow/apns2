package apns2_test

import (
	"fmt"
	apns "github.com/sideshow/apns2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func mockNotification() *apns.Notification {
	n := &apns.Notification{}
	n.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	n.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	return n
}

func mockClient(url string) *apns.Client {
	return &apns.Client{Host: url, HttpClient: http.DefaultClient}
}

func TestURL(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Incorrect Method", r.Method)
		}
		if r.URL.String() != fmt.Sprintf("/3/device/%s", n.DeviceToken) {
			t.Error("Incorrect URL", r.URL.String())
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
}

func TestDefaultHeaders(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Error("Header Content-Type should be application/json; charset=utf-8")
		}
		if r.Header.Get("apns-id") != "" {
			t.Error("Header apns-id should be unset")
		}
		if r.Header.Get("apns-priority") != "" {
			t.Error("Header apns-priority should be unset")
		}
		if r.Header.Get("apns-topic") != "" {
			t.Error("Header apns-topic should be unset")
		}
		if r.Header.Get("apns-expiration") != "" {
			t.Error("Header apns-expiration should be unset")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
}

func TestHeaders(t *testing.T) {
	n := mockNotification()
	n.ApnsId = "84DB694F-464F-49BD-960A-D6DB028335C9"
	n.Topic = "com.testapp"
	n.Priority = 10
	n.Expiration = time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("apns-id") != n.ApnsId {
			t.Error("Header apns-id should be ", n.ApnsId)
		}
		if r.Header.Get("apns-priority") != "10" {
			t.Error("Header apns-priority should be 10")
		}
		if r.Header.Get("apns-topic") != n.Topic {
			t.Error("Header apns-topic should be ", n.Topic)
		}
		if r.Header.Get("apns-expiration") != fmt.Sprintf("%v", n.Expiration.Unix()) {
			t.Error("Header apns-expiration should be ", n.Expiration.Unix())
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
}

func TestPayload(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(body, n.Payload) {
			t.Error("Body should be ", body, string(body))
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
}

func Test200SuccessResponse(t *testing.T) {
	n := mockNotification()
	var apnsId = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsId)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Error("StatusCode should be 200")
	}
	if res.ApnsId != apnsId {
		t.Error("ApnsID should be ", apnsId)
	}
	if !res.Sent() {
		t.Error("Success should be true")
	}
}

func Test400BadRequestPayloadEmptyResponse(t *testing.T) {
	n := mockNotification()
	var apnsId = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"reason\":\"PayloadEmpty\"}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 400 {
		t.Error("StatusCode should be 400")
	}
	if res.ApnsId != apnsId {
		t.Error("ApnsID should be ", apnsId)
	}
	if res.Reason != apns.ReasonPayloadEmpty {
		t.Error("Reason should be", apns.ReasonPayloadEmpty)
	}
	if res.Sent() {
		t.Error("Success should be false")
	}
}

func Test410UnregisteredResponse(t *testing.T) {
	n := mockNotification()
	var apnsId = "9F595474-356C-485E-B67F-9870BAE68702"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsId)
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("{\"reason\":\"Unregistered\", \"timestamp\":\"1421147681\"}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 410 {
		t.Error("StatusCode should be 410")
	}
	if res.ApnsId != apnsId {
		t.Error("ApnsID should be ", apnsId)
	}
	if res.Reason != apns.ReasonUnregistered {
		t.Error("Reason should be", apns.ReasonUnregistered)
	}
	if res.Timestamp.Unix() != 1421147681 {
		t.Error("Timestamp should be", 1421147681)
	}
	if res.Sent() {
		t.Error("Success should be false")
	}
}
