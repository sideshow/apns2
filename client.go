package apns2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	HostDevelopment = "https://api.development.push.apple.com"
	HostProduction  = "https://api.push.apple.com"
)

type Client struct {
	HttpClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

func NewClient(certificate tls.Certificate) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &Client{
		HttpClient:  &http.Client{Transport: transport},
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
	r.Header.Set("Content-Type", "application/json")
	if n.Topic != "" {
		r.Header.Set("apns-topic", n.Topic)
	}
	if n.Id != "" {
		r.Header.Set("apns-id", n.Id)
	}
	if n.Priority > 0 {
		r.Header.Set("apns-priority", fmt.Sprintf("%v", n.Priority))
	}
}

func (c *Client) Push(n *Notification) (*Response, error) {
	url := fmt.Sprintf("%v/3/device/%v", c.Host, n.DeviceToken)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(n.Payload))
	setHeaders(req, n)
	httpRes, httpErr := c.HttpClient.Do(req)

	if httpErr != nil {
		return nil, httpErr
	}
	defer httpRes.Body.Close()

	res := &Response{}
	res.StatusCode = httpRes.StatusCode
	res.NotificationID = httpRes.Header.Get("apns-id")
	if res.StatusCode == http.StatusOK {
		return res, nil
	} else {
		err := &APNSError{}
		json.NewDecoder(httpRes.Body).Decode(err)
		return res, err
	}
}
