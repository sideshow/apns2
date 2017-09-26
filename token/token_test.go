package token_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sideshow/apns2/token"
	"github.com/stretchr/testify/assert"
)

// AuthToken

func TestValidTokenFromP8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-valid.p8")
	assert.NoError(t, err)
}

func TestValidTokenFromP8Bytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/authkey-valid.p8")
	_, err := token.AuthKeyFromBytes(bytes)
	assert.NoError(t, err)
}

func TestNoSuchFileP8File(t *testing.T) {
	token, err := token.AuthKeyFromFile("")
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
	assert.Nil(t, token)
}

func TestInvalidP8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-invalid.p8")
	assert.Error(t, err)
}

func TestInvalidPKCS8P8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-invalid-pkcs8.p8")
	assert.Error(t, err)
}

func TestInvalidECDSAP8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-invalid-ecdsa.p8")
	assert.Error(t, err)
}

// Expiry & Generation

func TestExpired(t *testing.T) {
	token := &token.Token{}
	assert.True(t, token.Expired())
}

func TestNotExpired(t *testing.T) {
	token := &token.Token{
		IssuedAt: time.Now().Unix(),
	}
	assert.False(t, token.Expired())
}

func TestExpiresBeforeAnHour(t *testing.T) {
	token := &token.Token{
		IssuedAt: time.Now().Add(-50 * time.Minute).Unix(),
	}
	assert.True(t, token.Expired())
}

func TestGenerateIfExpired(t *testing.T) {
	authKey, _ := token.AuthKeyFromFile("_fixtures/authkey-valid.p8")
	token := &token.Token{
		AuthKey: authKey,
	}
	token.GenerateIfExpired()
	assert.Equal(t, time.Now().Unix(), token.IssuedAt)
}

func TestGenerateWithNoAuthKey(t *testing.T) {
	token := &token.Token{}
	bool, err := token.Generate()
	assert.False(t, bool)
	assert.Error(t, err)
}

func TestGenerateWithInvalidAuthKey(t *testing.T) {
	pubkeyCurve := elliptic.P521()
	privatekey := &ecdsa.PrivateKey{}
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)

	token := &token.Token{
		AuthKey: privatekey,
	}
	bool, err := token.Generate()
	assert.False(t, bool)
	assert.Error(t, err)
}
