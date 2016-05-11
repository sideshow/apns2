package apns2

import (
	"container/list"
	"crypto/sha1"
	"crypto/tls"
	"sync"
	"time"
)

type poolItem struct {
	key      [sha1.Size]byte
	client   *Client
	lastUsed time.Time
}

// ClientPool is a way to manage multiple connections to the APNs.
type ClientPool struct {
	// MaxSize is the maximum number of clients allowed in the pool. When this
	// limit is reached, the least recently used client is evicted. Set zero
	// for no limit.
	MaxSize int

	// MaxAge is the maximum age of clients in the pool. Upon retrieval, if a
	// client has remained unused in the pool for this duration or longer, it
	// is evicted and nil is returned. Set zero to disable this functionality.
	MaxAge time.Duration

	// Factory is the function which constructs clients if not found in the
	// pool.
	Factory func(certificate tls.Certificate) *Client

	cache map[[sha1.Size]byte]*list.Element
	ll    *list.List
	m     sync.Mutex
}

// NewClientPool returns a new ClientPool for prolonged, concurrent usage of
// multiple APNs clients. ClientPool is flexible enough to work best for your
// use case. When a client is not found in the pool, Get will return the result
// of calling Factory, which can be a Client or nil.
//
// Having multiple clients per certificate in the pool is not allowed.
//
// By default, MaxSize is 64, MaxAge is 10 minutes, and Factory always returns
// a Client with default options.
func NewClientPool() *ClientPool {
	pool := &ClientPool{
		MaxSize: 64,
		MaxAge:  10 * time.Minute,
		Factory: NewClient,
	}

	pool.initInternals()

	return pool
}

// Add adds a Client to the pool. You can use this to individually configure
// Clients in the pool.
func (p *ClientPool) Add(client *Client) {
	if p.cache == nil {
		p.initInternals()
	}
	p.m.Lock()
	defer p.m.Unlock()
	key := cacheKey(client.Certificate)
	now := time.Now()
	if ele, hit := p.cache[key]; hit {
		item := ele.Value.(*poolItem)
		item.client = client
		item.lastUsed = now
		p.ll.MoveToFront(ele)
		return
	}
	ele := p.ll.PushFront(&poolItem{key, client, now})
	p.cache[key] = ele
	if p.MaxSize != 0 && p.ll.Len() > p.MaxSize {
		p.m.Unlock()
		p.removeOldest()
		p.m.Lock()
	}
}

// Get gets a Client from the pool. If a Client is not found in the pool or if
// a Client has remained in the pool longer than MaxAge, Get will call the
// ClientPool's Factory function, store the result in the pool if non-nil, and
// return it.
func (p *ClientPool) Get(certificate tls.Certificate) *Client {
	if p.cache == nil {
		p.initInternals()
	}
	p.m.Lock()
	defer p.m.Unlock()
	key := cacheKey(certificate)
	now := time.Now()
	if ele, hit := p.cache[key]; hit {
		item := ele.Value.(*poolItem)
		if p.MaxAge != 0 && item.lastUsed.Before(now.Add(-p.MaxAge)) {
			c := p.Factory(certificate)
			if c == nil {
				return nil
			}
			item.client = c
		}
		item.lastUsed = now
		p.ll.MoveToFront(ele)
		return item.client
	}

	c := p.Factory(certificate)
	if c == nil {
		return nil
	}
	p.m.Unlock()
	p.Add(c)
	p.m.Lock()
	return c
}

// Len returns the current size of the ClientPool.
func (p *ClientPool) Len() int {
	if p.cache == nil {
		return 0
	}
	p.m.Lock()
	defer p.m.Unlock()
	return p.ll.Len()
}

func (p *ClientPool) initInternals() {
	p.cache = map[[sha1.Size]byte]*list.Element{}
	p.ll = list.New()
	p.m = sync.Mutex{}
}

func (p *ClientPool) removeOldest() {
	p.m.Lock()
	ele := p.ll.Back()
	p.m.Unlock()
	if ele != nil {
		p.removeElement(ele)
	}
}

func (p *ClientPool) removeElement(e *list.Element) {
	p.m.Lock()
	defer p.m.Unlock()
	p.ll.Remove(e)
	delete(p.cache, e.Value.(*poolItem).key)
}

func cacheKey(certificate tls.Certificate) [sha1.Size]byte {
	var data []byte

	for _, cert := range certificate.Certificate {
		data = append(data, cert...)
	}

	return sha1.Sum(data)
}
