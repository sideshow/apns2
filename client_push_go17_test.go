// +build go1.7

package apns2_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_PushWithCtx_WithTimeout(t *testing.T) {
	const timeout = time.Nanosecond
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	time.Sleep(timeout)
	res, err := mockClient(server.URL).PushWithCtx(n, ctx)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestClient_PushWithCtx(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	res, err := mockClient(server.URL).PushWithCtx(n, context.Background())
	assert.Nil(t, err)
	assert.Equal(t, res.ApnsID, apnsID)
}
