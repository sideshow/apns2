package apns2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

const (
	HostDevelopment = "https://api.development.push.apple.com"
	HostProduction  = "https://api.push.apple.com"
)

type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

func NewClient(certificate tls.Certificate) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	return &Client{
		HTTPClient:  &http.Client{Transport: transport},
		Certificate: certificate,
		Host:        HostDevelopment,
	}
}

func (c *Client) Development() *Client {
	c.Host = HostDevelopment
	return c
}

func (c *Client) Production() *Client {
	c.Host = HostProduction
	return c
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

func (c *Client) Push(n *Notification) (*Response, error) {
	url := fmt.Sprintf("%v/3/device/%v", c.Host, n.DeviceToken)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(n.Payload))
	setHeaders(req, n)
	httpRes, httpErr := c.HTTPClient.Do(req)

	if httpErr != nil {
		return nil, httpErr
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
