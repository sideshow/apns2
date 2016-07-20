// +build appengine

package apns2

import (
	"crypto/tls"
	"net"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/appengine/socket"
)

// GAEClient embeds an apns2.Client to use the special socket.Dial() method to connect to Apple
// Arbitratry net.Conn's are not allowed in the Google App Engine Classic environemtn
// https://cloud.google.com/appengine/docs/go/sockets
type GAEClient struct {
	*Client
	GConn *socket.Conn
	Ctx   context.Context
}

// SetContext assigns a new context to the underlying socket.Conn
//  client := NewGAEClient(cert)
//  client.SetContext(ctx)
//  client.Push(notification)
func (gclient *GAEClient) SetContext(ctx context.Context) {
	gclient.Ctx = ctx
	if gclient.GConn != nil {
		gclient.GConn.SetContext(gclient.Ctx)
	}
}

// NewGAEClient returns a new GAEClient with an underlying http.Client configured with
// the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
func NewGAEClient(certificate tls.Certificate) *GAEClient {
	gclient := &GAEClient{Client: NewClient(certificate)}

	transport := gclient.Client.HTTPClient.Transport.(*http2.Transport)
	transport.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
		gConn, err := socket.Dial(gclient.Ctx, network, addr)
		if err != nil {
			return nil, err
		}
		gclient.GConn = gConn

		tlsConn := tls.Client(gConn, cfg)
		return tlsConn, nil
	}

	return gclient
}

// Development sets the Client to use the APNs development push endpoint.
func (c *GAEClient) Development() *GAEClient {
	c.Host = HostDevelopment
	return c
}

// Production sets the Client to use the APNs production push endpoint.
func (c *GAEClient) Production() *GAEClient {
	c.Host = HostProduction
	return c
}
