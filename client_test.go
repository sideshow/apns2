package apns2_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/http2"

	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/stretchr/testify/assert"
)

// Mocks

func mockNotification() *apns.Notification {
	n := &apns.Notification{}
	n.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	n.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	return n
}

func mockCert() tls.Certificate {
	return tls.Certificate{}
}

func mockClient(url string) *apns.Client {
	return &apns.Client{Host: url, HTTPClient: http.DefaultClient}
}

// Unit Tests

func TestClientDefaultHost(t *testing.T) {
	client := apns.NewClient(mockCert())
	assert.Equal(t, "https://api.development.push.apple.com", client.Host)
}

func TestClientDevelopmentHost(t *testing.T) {
	client := apns.NewClient(mockCert()).Development()
	assert.Equal(t, "https://api.development.push.apple.com", client.Host)
}

func TestClientProductionHost(t *testing.T) {
	client := apns.NewClient(mockCert()).Production()
	assert.Equal(t, "https://api.push.apple.com", client.Host)
}

func TestClientBadUrlError(t *testing.T) {
	n := mockNotification()
	res, err := mockClient("badurl://badurl.com").Push(n)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestClientBadTransportError(t *testing.T) {
	n := mockNotification()
	client := mockClient("badurl://badurl.com")
	client.HTTPClient.Transport = nil
	res, err := client.Push(n)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestClientNameToCertificate(t *testing.T) {
	certificate, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	client := apns.NewClient(certificate)
	name := client.HTTPClient.Transport.(*http2.Transport).TLSClientConfig.NameToCertificate
	assert.Len(t, name, 1)

	certificate2 := tls.Certificate{}
	client2 := apns.NewClient(certificate2)
	name2 := client2.HTTPClient.Transport.(*http2.Transport).TLSClientConfig.NameToCertificate
	assert.Len(t, name2, 0)
}

// Functional Tests

func TestURL(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, fmt.Sprintf("/3/device/%s", n.DeviceToken), r.URL.String())
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestDefaultHeaders(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json; charset=utf-8", r.Header.Get("Content-Type"))
		assert.Equal(t, "", r.Header.Get("apns-id"))
		assert.Equal(t, "", r.Header.Get("apns-priority"))
		assert.Equal(t, "", r.Header.Get("apns-topic"))
		assert.Equal(t, "", r.Header.Get("apns-expiration"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestHeaders(t *testing.T) {
	n := mockNotification()
	n.ApnsID = "84DB694F-464F-49BD-960A-D6DB028335C9"
	n.Topic = "com.testapp"
	n.Priority = 10
	n.Expiration = time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, n.ApnsID, r.Header.Get("apns-id"))
		assert.Equal(t, "10", r.Header.Get("apns-priority"))
		assert.Equal(t, n.Topic, r.Header.Get("apns-topic"))
		assert.Equal(t, fmt.Sprintf("%v", n.Expiration.Unix()), r.Header.Get("apns-expiration"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPayload(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, n.Payload, body)
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestBadPayload(t *testing.T) {
	n := mockNotification()
	n.Payload = func() {}
	_, err := mockClient("").Push(n)
	assert.Error(t, err)
}

func Test200SuccessResponse(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, apnsID, res.ApnsID)
	assert.Equal(t, true, res.Sent())
}

func Test400BadRequestPayloadEmptyResponse(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"reason\":\"PayloadEmpty\"}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, apnsID, res.ApnsID)
	assert.Equal(t, apns.ReasonPayloadEmpty, res.Reason)
	assert.Equal(t, false, res.Sent())
}

func Test410UnregisteredResponse(t *testing.T) {
	n := mockNotification()
	var apnsID = "9F595474-356C-485E-B67F-9870BAE68702"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("{\"reason\":\"Unregistered\", \"timestamp\": 1458114061260 }"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
	assert.Equal(t, 410, res.StatusCode)
	assert.Equal(t, apnsID, res.ApnsID)
	assert.Equal(t, apns.ReasonUnregistered, res.Reason)
	assert.Equal(t, int64(1458114061260)/1000, res.Timestamp.Unix())
	assert.Equal(t, false, res.Sent())
}

func TestMalformedJSONResponse(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte("{{MalformedJSON}}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	assert.Error(t, err)
	assert.Equal(t, false, res.Sent())
}
