// Package apns2 is a go Apple Push Notification Service (APNs) provider that
// allows you to send remote notifications to your iOS, tvOS, and OS X
// apps, using the new APNs HTTP/2 network protocol.
package apns2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// Apple HTTP/2 Development & Production urls
const (
	HostDevelopment = "https://api.development.push.apple.com"
	HostProduction  = "https://api.push.apple.com"
)

// DefaultHost is a mutable var for testing purposes
var DefaultHost = HostDevelopment

// Client represents a connection with the APNs
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

// NewClient returns a new Client with an underlying http.Client configured with
// the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
func NewClient(certificate tls.Certificate) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return tls.DialWithDialer(&net.Dialer{Timeout: 1 * time.Second}, network, addr, cfg)
		},
	}
	return &Client{
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   1 * time.Second,
		},
		Certificate: certificate,
		Host:        DefaultHost,
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
func (c *Client) Push(n *Notification) (*Response, error) {
	payload, err := json.Marshal(n)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v/3/device/%v", c.Host, n.DeviceToken)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	setHeaders(req, n)
	httpRes, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	response := &Response{}
	response.StatusCode = httpRes.StatusCode
	response.ApnsID = httpRes.Header.Get("apns-id")

	decoder := json.NewDecoder(httpRes.Body)
	if err := decoder.Decode(&response); err != nil && err != io.EOF {
		return &Response{}, err
	}
	return response, nil
}

func setHeaders(r *http.Request, n *Notification) {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	if n.Topic != "" {
		r.Header.Set("apns-topic", n.Topic)
	}
	if n.ApnsID != "" {
		r.Header.Set("apns-id", n.ApnsID)
	}
	if n.Priority > 0 {
		r.Header.Set("apns-priority", fmt.Sprintf("%v", n.Priority))
	}
	if !n.Expiration.IsZero() {
		r.Header.Set("apns-expiration", fmt.Sprintf("%v", n.Expiration.Unix()))
	}
}
