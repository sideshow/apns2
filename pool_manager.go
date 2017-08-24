package apns2

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/net/http2"
	"golang.org/x/net/idna"
)

type poolManager struct {
	connPool http2.ClientConnPool
	ctx      context.Context
	poolMu   *sync.Mutex
	u        *url.URL
}

// directly copied from source
func authorityAddr(scheme string, authority string) (addr string) {
	host, port, err := net.SplitHostPort(authority)
	if err != nil {
		port = "443"
		if scheme == "http" {
			port = "80"
		}
		host = authority
	}
	if a, err := idna.ToASCII(host); err == nil {
		host = a
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host + ":" + port
	}
	return net.JoinHostPort(host, port)
}

func newPoolManager(transport *http2.Transport, environment string) (*poolManager, error) {
	transport.CloseIdleConnections()
	rf := reflect.Indirect(reflect.ValueOf(transport)).FieldByName("connPoolOrDef")
	connPool := *(*http2.ClientConnPool)(unsafe.Pointer(rf.UnsafeAddr()))
	rf = reflect.Indirect(reflect.ValueOf(connPool)).FieldByName("mu")
	poolMu := (*sync.Mutex)(unsafe.Pointer(rf.UnsafeAddr()))
	u, err := url.Parse(environment)
	if err != nil {
		return nil, err
	}
	return &poolManager{
		connPool: connPool,
		u:        u,
		poolMu:   poolMu,
		ctx:      context.Background(),
	}, nil
}

func (pm *poolManager) addNewConn() error {
	rff := reflect.Indirect(reflect.ValueOf(pm.connPool)).FieldByName("conns")
	pm.poolMu.Lock()
	internalConns := *(*map[string][]*http2.ClientConn)(unsafe.Pointer(rff.UnsafeAddr()))
	pm.poolMu.Unlock()
	for _, conns := range internalConns {
		for _, conn := range conns {
			rv := reflect.Indirect(reflect.ValueOf(conn))
			rf := rv.FieldByName("mu")
			mu := (*sync.Mutex)(unsafe.Pointer(rf.UnsafeAddr()))
			rf = rv.FieldByName("closed")
			mu.Lock()
			rf = reflect.Indirect(reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())))
			rf.SetBool(true)
			mu.Unlock()
			defer func() {
				mu.Lock()
				rf.SetBool(false)
				mu.Unlock()
			}()
		}
	}
	cc, err := pm.connPool.GetClientConn(&http.Request{Close: false}, authorityAddr(pm.u.Scheme, pm.u.Host))
	if err != nil {
		return err
	}
	go pm.pingConn(cc)
	return nil
}

func (pm *poolManager) pingConn(cc *http2.ClientConn) {
	for {
		err := cc.Ping(pm.ctx)
		if err != nil {
			pm.addNewConn()
			return
		}
		time.Sleep(PingInverval)
	}
}
