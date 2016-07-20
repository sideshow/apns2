// +build appengine

package apns2_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"google.golang.org/appengine/aetest"
)

// Mocks

func mockGAEClient(url string) *apns.GAEClient {
	cert, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	gclient := apns.NewGAEClient(cert)
	gclient.Client.Host = url
	return gclient
}

// Unit Tests

func TestGAEClientDefaultHost(t *testing.T) {
	client := apns.NewGAEClient(mockCert())
	assert.Equal(t, "https://api.development.push.apple.com", client.Host)
}

func TestGAEClientDevelopmentHost(t *testing.T) {
	client := apns.NewGAEClient(mockCert()).Development()
	assert.Equal(t, "https://api.development.push.apple.com", client.Host)
}

func TestGAEClientProductionHost(t *testing.T) {
	client := apns.NewGAEClient(mockCert()).Production()
	assert.Equal(t, "https://api.push.apple.com", client.Host)
}

func TestGAEConnection(t *testing.T) {
	n := mockNotification()

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("apns-id", "XXXXXXXXXXXXXXX")
	}))
	server.TLS = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "h2-14"},
	}
	server.StartTLS()
	defer server.Close()

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	gclient := mockGAEClient(server.URL)

	transport := gclient.Client.HTTPClient.Transport.(*http2.Transport)
	transport.TLSClientConfig.InsecureSkipVerify = true
	transport.TLSClientConfig.NextProtos = []string{"h2", "h2-14"}

	gclient.SetContext(ctx)
	_, err = gclient.Push(n)
	assert.NoError(t, err)
}
