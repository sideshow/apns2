package apns2

import (
	"container/list"
	"crypto/sha1"
	"crypto/tls"
	"sync"
	"time"
)

type managerItem struct {
	key      [sha1.Size]byte
	client   *Client
	lastUsed time.Time
}

// ClientFactory is a the func that provides to get a new client
type ClientFactory func(tls.Certificate) *Client

// ManagerOpts is the func to set options to the connection manager
type ManagerOpts func(c *manager) error

// MaxSize is the maximum number of clients allowed in the manager. When
// this limit is reached, the least recently used client is evicted. Set
// zero for no limit.
func MaxSize(size int) ManagerOpts {
	return func(c *manager) error {
		c.maxSize = size
		return nil
	}
}

// MaxAge is the maximum age of clients in the manager. Upon retrieval, if
// a client has remained unused in the manager for this duration or longer,
// it is evicted and nil is returned. Set zero to disable this
// functionality.
func MaxAge(age time.Duration) ManagerOpts {
	return func(c *manager) error {
		c.maxAge = age
		return nil
	}
}

// Factory is the function which constructs clients if not found in the
// manager.
func Factory(f ClientFactory) ManagerOpts {
	return func(c *manager) error {
		c.factory = f
		return nil
	}
}

// ClientManager interface provides work with multiple connection to the APNS
type ClientManager interface {
	Add(*Client)
	Get(tls.Certificate) *Client
	Len() int
}

// manager is a way to manage multiple connections to the APNs.
type manager struct {
	// maxSize is the maximum number of clients allowed in the manager. When
	// this limit is reached, the least recently used client is evicted. Set
	// zero for no limit.
	maxSize int

	// maxAge is the maximum age of clients in the manager. Upon retrieval, if
	// a client has remained unused in the manager for this duration or longer,
	// it is evicted and nil is returned. Set zero to disable this
	// functionality.
	maxAge time.Duration

	// factory is the function which constructs clients if not found in the
	// manager.
	factory ClientFactory

	ll *list.List
	// mutex
	mu    sync.Mutex
	cache map[[sha1.Size]byte]*list.Element
}

// NewClientManager returns a new ClientManager for prolonged, concurrent usage
// of multiple APNs clients. ClientManager is flexible enough to work best for
// your use case. When a client is not found in the manager, Get will return
// the result of calling Factory, which can be a Client or nil.
//
// Having multiple clients per certificate in the manager is not allowed.
//
// By default, MaxSize is 64, MaxAge is 10 minutes, and Factory always returns
// a Client with default options.
//
// You can to pass a manager opts, that sets up the size, age and factory
func NewClientManager(opts ...ManagerOpts) ClientManager {
	manager := &manager{
		maxSize: 64,
		maxAge:  10 * time.Minute,
		factory: NewClient,
		cache:   make(map[[sha1.Size]byte]*list.Element),
		ll:      list.New(),
	}

	// apply options to manager
	for _, o := range opts {
		o(manager)
	}

	return manager
}

// Add adds a Client to the manager. You can use this to individually configure
// Clients in the manager.
func (m *manager) Add(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := cacheKey(client.Certificate)
	now := time.Now()
	if ele, hit := m.cache[key]; hit {
		item := ele.Value.(*managerItem)
		item.client = client
		item.lastUsed = now
		m.ll.MoveToFront(ele)
		return
	}
	ele := m.ll.PushFront(&managerItem{key, client, now})
	m.cache[key] = ele
	if m.maxSize != 0 && m.ll.Len() > m.maxSize {
		m.mu.Unlock()
		m.removeOldest()
		m.mu.Lock()
	}
}

// Get gets a Client from the manager. If a Client is not found in the manager
// or if a Client has remained in the manager longer than MaxAge, Get will call
// the ClientManager's Factory function, store the result in the manager if
// non-nil, and return it.
func (m *manager) Get(certificate tls.Certificate) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := cacheKey(certificate)
	now := time.Now()
	if ele, hit := m.cache[key]; hit {
		item := ele.Value.(*managerItem)
		if m.maxAge != 0 && item.lastUsed.Before(now.Add(-m.maxAge)) {
			c := m.factory(certificate)
			if c == nil {
				return nil
			}
			item.client = c
		}
		item.lastUsed = now
		m.ll.MoveToFront(ele)
		return item.client
	}

	c := m.factory(certificate)
	if c == nil {
		return nil
	}
	m.mu.Unlock()
	m.Add(c)
	m.mu.Lock()
	return c
}

// Len returns the current size of the ClientManager.
func (m *manager) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ll.Len()
}

func (m *manager) removeOldest() {
	m.mu.Lock()
	ele := m.ll.Back()
	m.mu.Unlock()
	if ele != nil {
		m.removeElement(ele)
	}
}

func (m *manager) removeElement(e *list.Element) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ll.Remove(e)
	delete(m.cache, e.Value.(*managerItem).key)
}

func cacheKey(certificate tls.Certificate) [sha1.Size]byte {
	var data []byte

	for _, cert := range certificate.Certificate {
		data = append(data, cert...)
	}

	return sha1.Sum(data)
}
