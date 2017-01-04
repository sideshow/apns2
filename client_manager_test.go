package apns2_test

import (
	"bytes"
	"crypto/tls"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/stretchr/testify/assert"
)

func TestNewClientManager(t *testing.T) {
	manager := apns2.NewClientManager()
	assert.Equal(t, manager.MaxSize, 64)
	assert.Equal(t, manager.MaxAge, 10*time.Minute)
}

func TestClientManagerGetWithoutNew(t *testing.T) {
	manager := apns2.ClientManager{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	c1 := manager.Get(mockCert())
	c2 := manager.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	assert.NotNil(t, c1)
	assert.Equal(t, v1.Pointer(), v2.Pointer())
	assert.Equal(t, 1, manager.Len())
}

func TestClientManagerAddWithoutNew(t *testing.T) {
	wg := sync.WaitGroup{}

	manager := apns2.ClientManager{
		MaxSize: 1,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			manager.Add(apns2.NewClient(mockCert()))
			assert.Equal(t, 1, manager.Len())
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestClientManagerLenWithoutNew(t *testing.T) {
	manager := apns2.ClientManager{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	assert.Equal(t, 0, manager.Len())
}

func TestClientManagerGetDefaultOptions(t *testing.T) {
	manager := apns2.NewClientManager()
	c1 := manager.Get(mockCert())
	c2 := manager.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	assert.NotNil(t, c1)
	assert.Equal(t, v1.Pointer(), v2.Pointer())
	assert.Equal(t, 1, manager.Len())
}

func TestClientManagerGetNilClientFactory(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.Factory = func(certificate tls.Certificate) *apns2.Client {
		return nil
	}
	c1 := manager.Get(mockCert())
	c2 := manager.Get(mockCert())
	assert.Nil(t, c1)
	assert.Nil(t, c2)
	assert.Equal(t, 0, manager.Len())
}

func TestClientManagerGetMaxAgeExpiration(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.MaxAge = time.Nanosecond
	c1 := manager.Get(mockCert())
	time.Sleep(time.Microsecond)
	c2 := manager.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	assert.NotNil(t, c1)
	assert.NotEqual(t, v1.Pointer(), v2.Pointer())
	assert.Equal(t, 1, manager.Len())
}

func TestClientManagerGetMaxAgeExpirationWithNilFactory(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.Factory = func(certificate tls.Certificate) *apns2.Client {
		return nil
	}
	manager.MaxAge = time.Nanosecond
	manager.Add(apns2.NewClient(mockCert()))
	c1 := manager.Get(mockCert())
	time.Sleep(time.Microsecond)
	c2 := manager.Get(mockCert())
	assert.Nil(t, c1)
	assert.Nil(t, c2)
	assert.Equal(t, 1, manager.Len())
}

func TestClientManagerGetMaxSizeExceeded(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.MaxSize = 1
	cert1 := mockCert()
	_ = manager.Get(cert1)
	cert2, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	_ = manager.Get(cert2)
	cert3, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid-encrypted.p12", "password")
	c := manager.Get(cert3)
	assert.True(t, bytes.Equal(cert3.Certificate[0], c.Certificate.Certificate[0]))
	assert.Equal(t, 1, manager.Len())
}

func TestClientManagerAdd(t *testing.T) {
	fn := func(certificate tls.Certificate) *apns2.Client {
		t.Fatal("factory should not have been called")
		return nil
	}

	manager := apns2.NewClientManager()
	manager.Factory = fn
	manager.Add(apns2.NewClient(mockCert()))
	manager.Get(mockCert())
}

func TestClientManagerAddTwice(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.Add(apns2.NewClient(mockCert()))
	manager.Add(apns2.NewClient(mockCert()))
	assert.Equal(t, 1, manager.Len())
}
