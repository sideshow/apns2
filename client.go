// Package apns2 is a go Apple Push Notification Service (APNs) provider that
// allows you to send remote notifications to your iOS, tvOS, and OS X
// apps, using the new APNs HTTP/2 network protocol.
package apns2

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ringsaturn/apns2/token"
)

// Apple HTTP/2 Development & Production urls
const (
	HostDevelopment = "https://api.sandbox.push.apple.com"
	HostProduction  = "https://api.push.apple.com"
)

// DefaultHost is a mutable var for testing purposes
var DefaultHost = HostDevelopment

var (
	// HTTPClientTimeout specifies a time limit for requests made by the
	// HTTPClient. The timeout includes connection time, any redirects,
	// and reading the response body.
	HTTPClientTimeout = 60 * time.Second

	// ReadIdleTimeout is the timeout after which a health check using a ping
	// frame will be carried out if no frame is received on the connection. If
	// zero, no health check is performed.
	ReadIdleTimeout = 15 * time.Second

	// TCPKeepAlive specifies the keep-alive period for an active network
	// connection. If zero, keep-alive probes are sent with a default value
	// (currently 15 seconds)
	TCPKeepAlive = 15 * time.Second

	// TLSDialTimeout is the maximum amount of time a dial will wait for a connect
	// to complete.
	TLSDialTimeout = 20 * time.Second
)

// Client represents a connection with the APNs
type Client struct {
	Host        string
	Certificate tls.Certificate
	Token       *token.Token
	HTTPClient  *http.Client
}

type connectionCloser interface {
	CloseIdleConnections()
}

// NewClient returns a new Client with an underlying http.Client configured with
// the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
//
// If your use case involves multiple long-lived connections, consider using
// the ClientManager, which manages clients for you.
func NewClient(certificate tls.Certificate) *Client {
	client := &Client{
		HTTPClient: &http.Client{
			Timeout:   HTTPClientTimeout,
			Transport: &http.Transport{},
		},
		Certificate: certificate,
		Host:        DefaultHost,
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	client.HTTPClient.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	return client
}

// NewTokenClient returns a new Client with an underlying http.Client configured
// with the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
func NewTokenClient(token *token.Token) *Client {
	return &Client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: HTTPClientTimeout,
		},
		Host: DefaultHost,
	}
}

// Development sets the Client to use the APNs development push endpoint.
func (c *Client) Development() *Client {
	c.Host = HostDevelopment
	return c
}

// Production sets the Client to use the APNs production push endpoint.
func (c *Client) Production() *Client {
	c.Host = HostProduction
	return c
}

// Push sends a Notification to the APNs gateway. If the underlying http.Client
// is not currently connected, this method will attempt to reconnect
// transparently before sending the notification. It will return a Response
// indicating whether the notification was accepted or rejected by the APNs
// gateway, or an error if something goes wrong.
//
// Use PushWithContext if you need better cancellation and timeout control.
func (c *Client) Push(n *Notification) (*Response, error) {
	return c.PushWithContext(context.Background(), n)
}

// PushWithContext sends a Notification to the APNs gateway. Context carries a
// deadline and a cancellation signal and allows you to close long running
// requests when the context timeout is exceeded. Context can be nil, for
// backwards compatibility.
//
// If the underlying http.Client is not currently connected, this method will
// attempt to reconnect transparently before sending the notification. It will
// return a Response indicating whether the notification was accepted or
// rejected by the APNs gateway, or an error if something goes wrong.
func (c *Client) PushWithContext(ctx context.Context, n *Notification) (*Response, error) {
	payload, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	url := c.Host + "/3/device/" + n.DeviceToken
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	if c.Token != nil {
		c.setTokenHeader(request)
	}

	setHeaders(request, n)

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	r := &Response{}
	r.StatusCode = response.StatusCode
	r.ApnsID = response.Header.Get("apns-id")
	r.ApnsUniqueID = response.Header.Get("apns-unique-id")

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(r); err != nil && err != io.EOF {
		return &Response{}, err
	}
	return r, nil
}

// CloseIdleConnections closes any underlying connections which were previously
// connected from previous requests but are now sitting idle. It will not
// interrupt any connections currently in use.
func (c *Client) CloseIdleConnections() {
	c.HTTPClient.Transport.(connectionCloser).CloseIdleConnections()
}

func (c *Client) setTokenHeader(r *http.Request) {
	bearer := c.Token.GenerateIfExpired()
	r.Header.Set("authorization", "bearer "+bearer)
}

func setHeaders(r *http.Request, n *Notification) {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	if n.Topic != "" {
		r.Header.Set("apns-topic", n.Topic)
	}
	if n.ApnsID != "" {
		r.Header.Set("apns-id", n.ApnsID)
	}
	if n.CollapseID != "" {
		r.Header.Set("apns-collapse-id", n.CollapseID)
	}
	if n.Priority > 0 {
		r.Header.Set("apns-priority", strconv.Itoa(n.Priority))
	}
	if !n.Expiration.IsZero() {
		r.Header.Set("apns-expiration", strconv.FormatInt(n.Expiration.Unix(), 10))
	}
	if n.PushType != "" {
		r.Header.Set("apns-push-type", string(n.PushType))
	} else {
		r.Header.Set("apns-push-type", string(PushTypeAlert))
	}

}
