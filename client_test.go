package apns2_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/http2"

	apns "github.com/ringsaturn/apns2"
	"github.com/ringsaturn/apns2/token"
	"github.com/stretchr/testify/assert"
)

// Mocks

func mockNotification() *apns.Notification {
	n := &apns.Notification{}
	n.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	n.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	return n
}

func mockToken() *token.Token {
	pubkeyCurve := elliptic.P256()
	authKey, _ := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	return &token.Token{AuthKey: authKey}
}

func mockCert() tls.Certificate {
	return tls.Certificate{}
}

func mockClient(url string) *apns.Client {
	return &apns.Client{Host: url, HTTPClient: http.DefaultClient}
}

type mockTransport struct {
	*http2.Transport
	closed bool
}

func (c *mockTransport) CloseIdleConnections() {
	c.closed = true
}

// Unit Tests

func TestClientDefaultHost(t *testing.T) {
	client := apns.NewClient(mockCert())
	assert.Equal(t, "https://api.sandbox.push.apple.com", client.Host)
}

func TestTokenDefaultHost(t *testing.T) {
	client := apns.NewTokenClient(mockToken()).Development()
	assert.Equal(t, "https://api.sandbox.push.apple.com", client.Host)
}

func TestClientDevelopmentHost(t *testing.T) {
	client := apns.NewClient(mockCert()).Development()
	assert.Equal(t, "https://api.sandbox.push.apple.com", client.Host)
}

func TestTokenClientDevelopmentHost(t *testing.T) {
	client := apns.NewTokenClient(mockToken()).Development()
	assert.Equal(t, "https://api.sandbox.push.apple.com", client.Host)
}

func TestClientProductionHost(t *testing.T) {
	client := apns.NewClient(mockCert()).Production()
	assert.Equal(t, "https://api.push.apple.com", client.Host)
}

func TestTokenClientProductionHost(t *testing.T) {
	client := apns.NewTokenClient(mockToken()).Production()
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

func TestClientBadDeviceToken(t *testing.T) {
	n := &apns.Notification{}
	n.DeviceToken = "DGw\aOoD+HwSroh#Ug]%xzd]"
	n.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	res, err := mockClient("https://api.push.apple.com").Push(n)
	assert.Error(t, err)
	assert.Nil(t, res)
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
		assert.Equal(t, "", r.Header.Get("apns-collapse-id"))
		assert.Equal(t, "", r.Header.Get("apns-priority"))
		assert.Equal(t, "", r.Header.Get("apns-topic"))
		assert.Equal(t, "", r.Header.Get("apns-expiration"))
		assert.Equal(t, "", r.Header.Get("thread-id"))
		assert.Equal(t, "alert", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestClientPushWithContextWithTimeout(t *testing.T) {
	const timeout = time.Nanosecond
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	time.Sleep(timeout)
	res, err := mockClient(server.URL).PushWithContext(ctx, n)
	assert.Error(t, err)
	assert.Nil(t, res)
	cancel()
}

func TestClientPushWithContext(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	res, err := mockClient(server.URL).PushWithContext(context.Background(), n)
	assert.Nil(t, err)
	assert.Equal(t, res.ApnsID, apnsID)
}

func TestClientPushWithNilContext(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	//lint:ignore SA1012 we need use nil context for test
	res, err := mockClient(server.URL).PushWithContext(nil, n)
	assert.EqualError(t, err, "net/http: nil Context")
	assert.Nil(t, res)
}

func TestHeaders(t *testing.T) {
	n := mockNotification()
	n.ApnsID = "84DB694F-464F-49BD-960A-D6DB028335C9"
	n.CollapseID = "game1.start.identifier"
	n.Topic = "com.testapp"
	n.Priority = 10
	n.Expiration = time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, n.ApnsID, r.Header.Get("apns-id"))
		assert.Equal(t, n.CollapseID, r.Header.Get("apns-collapse-id"))
		assert.Equal(t, "10", r.Header.Get("apns-priority"))
		assert.Equal(t, n.Topic, r.Header.Get("apns-topic"))
		assert.Equal(t, fmt.Sprintf("%v", n.Expiration.Unix()), r.Header.Get("apns-expiration"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeAlertHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeAlert
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "alert", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeBackgroundHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeBackground
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "background", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeLocationHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeLocation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "location", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeVOIPHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeVOIP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "voip", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeComplicationHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeComplication
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "complication", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeFileProviderHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeFileProvider
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "fileprovider", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeMDMHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeMDM
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "mdm", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestPushTypeLiveActivityHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeLiveActivity
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "liveactivity", r.Header.Get("apns-push-type"))
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	assert.NoError(t, err)
}

func TestAuthorizationHeader(t *testing.T) {
	n := mockNotification()
	token := mockToken()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json; charset=utf-8", r.Header.Get("Content-Type"))
		assert.Equal(t, fmt.Sprintf("bearer %v", token.Bearer), r.Header.Get("authorization"))
	}))
	defer server.Close()

	client := mockClient(server.URL)
	client.Token = token
	_, err := client.Push(n)
	assert.NoError(t, err)
}

func TestPayload(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
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
		_, _ = w.Write([]byte("{\"reason\":\"PayloadEmpty\"}"))
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
		_, _ = w.Write([]byte("{\"reason\":\"Unregistered\", \"timestamp\": 1458114061260 }"))
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
		_, _ = w.Write([]byte("{{MalformedJSON}}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	assert.Error(t, err)
	assert.Equal(t, false, res.Sent())
}

func TestCloseIdleConnections(t *testing.T) {
	transport := &mockTransport{}

	client := mockClient("")
	client.HTTPClient.Transport = transport

	assert.Equal(t, false, transport.closed)
	client.CloseIdleConnections()
	assert.Equal(t, true, transport.closed)
}
