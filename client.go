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
	"math/rand"
	"net"
	"net/http"
	"sync"
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

var (
	// TLSDialTimeout is the maximum amount of time a dial will wait for a connect
	// to complete.
	TLSDialTimeout = 20 * time.Second
	// HTTPClientTimeout specifies a time limit for requests made by the
	// HTTPClient. The timeout includes connection time, any redirects,
	// and reading the response body.
	HTTPClientTimeout = 30 * time.Second
	// PingPongFrequency is the interval with which a client will PING APNs
	// servers.
	PingPongFrequency = 15 * time.Second
)

// Client represents a connection with the APNs
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
	conn        net.Conn
	pinging     bool
	newConnChan chan struct{}
	stopChan    chan struct{}
	m           *sync.Mutex
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
//
// Alternatively, you can keep the clients connection healthy by calling
// EnablePinging, which will send PING frames to APNs servers with the interval
// specified via PingPongFrequency.
func NewClient(certificate tls.Certificate) (client *Client) {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	client = &Client{
		Certificate: certificate,
		Host:        DefaultHost,
		newConnChan: make(chan struct{}),
		stopChan:    make(chan struct{}),
		m:           new(sync.Mutex),
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS: func(network, addr string, cfg *tls.Config) (c net.Conn, e error) {
			c, e = tls.DialWithDialer(&net.Dialer{Timeout: TLSDialTimeout}, network, addr, cfg)
			if e == nil {
				client.m.Lock()
				defer client.m.Unlock()
				client.conn = c
				if client.pinging {
					client.newConnChan <- struct{}{}
				}
			}
			return
		},
	}
	client.HTTPClient = &http.Client{
		Transport: transport,
		Timeout:   HTTPClientTimeout,
	}
	return
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

// EnablePinging tries to send PING frames to APNs servers whenever the client
// has a valid connection. If the willHandleDrops parameter is set to true, this
// function returns a read-only channel that gets notified when pinging fails.
// This allows the user to take actions to preemptively reinitialize the client's
// connection. The second return value indicates whether the call has successfully
// enabled pinging.
func (c *Client) EnablePinging(willHandleDrops bool) (<-chan struct{}, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.pinging {
		return nil, false
	}
	c.pinging = true
	var dropSignal chan struct{}
	if willHandleDrops {
		dropSignal = make(chan struct{})
	}
	go func() {
		// 8 bytes of random data used for PING-PONG, as per HTTP/2 spec.
		data := [8]byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}
		pinger := new(time.Ticker)
		var framer *http2.Framer
		c.m.Lock()
		if c.conn != nil {
			framer = http2.NewFramer(c.conn, c.conn)
			pinger = time.NewTicker(PingPongFrequency)
		}
		c.m.Unlock()
		for {
			select {
			case <-pinger.C:
				err := framer.WritePing(false, data)
				if err != nil {
					// APNs did not answer with pong, which means the connection
					// has been dropped. Stop trying and notify the drop handler,
					// if there is any.
					c.m.Lock()
					c.conn = nil
					c.m.Unlock()
					pinger.Stop()
					if willHandleDrops {
						dropSignal <- struct{}{}
					}
				}
			case <-c.newConnChan:
				c.m.Lock()
				framer = http2.NewFramer(c.conn, c.conn)
				c.m.Unlock()
				pinger.Stop()
				pinger = time.NewTicker(PingPongFrequency)
			case <-c.stopChan:
				pinger.Stop()
				return
			}
		}
	}()
	return dropSignal, true
}

// DisablePinging stops the pinging operation associated with the client, if
// there's any, and returns a boolean that indicates if the call has successfully
// stopped the pinging operation.
func (c *Client) DisablePinging() bool {
	c.m.Lock()
	defer c.m.Unlock()
	if c.pinging {
		c.pinging = false
		c.stopChan <- struct{}{}
		return true
	}
	return false
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
	if n.ThreadID != "" {
		r.Header.Set("thread-id", n.ThreadID)
	}
	if n.Priority > 0 {
		r.Header.Set("apns-priority", fmt.Sprintf("%v", n.Priority))
	}
	if !n.Expiration.IsZero() {
		r.Header.Set("apns-expiration", fmt.Sprintf("%v", n.Expiration.Unix()))
	}
}
