package apns2_test

import (
	"bytes"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/stretchr/testify/assert"
)

func TestNewClientPool(t *testing.T) {
	pool := apns2.NewClientPool()
	assert.Equal(t, pool.MaxSize, 64)
	assert.Equal(t, pool.MaxAge, 10*time.Minute)
}

func TestClientPoolGetWithoutNew(t *testing.T) {
	pool := apns2.ClientPool{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	c1 := pool.Get(mockCert())
	c2 := pool.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	assert.NotNil(t, c1)
	assert.Equal(t, v1.Pointer(), v2.Pointer())
	assert.Equal(t, 1, pool.Len())
}

func TestClientPoolAddWithoutNew(t *testing.T) {
	pool := apns2.ClientPool{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	pool.Add(apns2.NewClient(mockCert()))
	assert.Equal(t, 1, pool.Len())
}

func TestClientPoolLenWithoutNew(t *testing.T) {
	pool := apns2.ClientPool{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	assert.Equal(t, 0, pool.Len())
}

func TestClientPoolGetDefaultOptions(t *testing.T) {
	pool := apns2.NewClientPool()
	c1 := pool.Get(mockCert())
	c2 := pool.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	assert.NotNil(t, c1)
	assert.Equal(t, v1.Pointer(), v2.Pointer())
	assert.Equal(t, 1, pool.Len())
}

func TestClientPoolGetNilClientFactory(t *testing.T) {
	pool := apns2.NewClientPool()
	pool.Factory = func(certificate tls.Certificate) *apns2.Client {
		return nil
	}
	c1 := pool.Get(mockCert())
	c2 := pool.Get(mockCert())
	assert.Nil(t, c1)
	assert.Nil(t, c2)
	assert.Equal(t, 0, pool.Len())
}

func TestClientPoolGetMaxAgeExpiration(t *testing.T) {
	pool := apns2.NewClientPool()
	pool.MaxAge = time.Nanosecond
	c1 := pool.Get(mockCert())
	time.Sleep(time.Microsecond)
	c2 := pool.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	assert.NotNil(t, c1)
	assert.NotEqual(t, v1.Pointer(), v2.Pointer())
	assert.Equal(t, 1, pool.Len())
}

func TestClientPoolGetMaxAgeExpirationWithNilFactory(t *testing.T) {
	pool := apns2.NewClientPool()
	pool.Factory = func(certificate tls.Certificate) *apns2.Client {
		return nil
	}
	pool.MaxAge = time.Nanosecond
	pool.Add(apns2.NewClient(mockCert()))
	c1 := pool.Get(mockCert())
	time.Sleep(time.Microsecond)
	c2 := pool.Get(mockCert())
	assert.Nil(t, c1)
	assert.Nil(t, c2)
	assert.Equal(t, 1, pool.Len())
}

func TestClientPoolGetMaxSizeExceeded(t *testing.T) {
	pool := apns2.NewClientPool()
	pool.MaxSize = 1
	cert1 := mockCert()
	_ = pool.Get(cert1)
	cert2, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	_ = pool.Get(cert2)
	cert3, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid-encrypted.p12", "password")
	c := pool.Get(cert3)
	assert.True(t, bytes.Equal(cert3.Certificate[0], c.Certificate.Certificate[0]))
	assert.Equal(t, 1, pool.Len())
}

func TestClientPoolAdd(t *testing.T) {
	fn := func(certificate tls.Certificate) *apns2.Client {
		t.Fatal("factory should not have been called")
		return nil
	}

	pool := apns2.NewClientPool()
	pool.Factory = fn
	pool.Add(apns2.NewClient(mockCert()))
	pool.Get(mockCert())
}

func TestClientPoolAddTwice(t *testing.T) {
	pool := apns2.NewClientPool()
	pool.Add(apns2.NewClient(mockCert()))
	pool.Add(apns2.NewClient(mockCert()))
	assert.Equal(t, 1, pool.Len())
}
