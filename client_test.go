package apns2_test

import (
	"fmt"
	apns "github.com/sideshow/apns2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Header Content-Type should be application/json")
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
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Error(err)
	}
}

func TestHeaders(t *testing.T) {
	n := mockNotification()
	n.Id = "84DB694F-464F-49BD-960A-D6DB028335C9"
	n.Topic = "com.testapp"
	n.Priority = 10
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("apns-id") != n.Id {
			t.Error("Header apns-id should be ", n.Id)
		}
		if r.Header.Get("apns-priority") != "10" {
			t.Error("Header apns-priority should be 10")
		}
		if r.Header.Get("apns-topic") != n.Topic {
			t.Error("Header apns-topic should be ", n.Topic)
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
