package token_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/sideshow/apns2/token"
)

// Expired is exported and documented, so callers can reasonably use it
// from monitoring goroutines while the client refreshes the token.
func TestTokenExpiredConcurrent(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tk := &token.Token{AuthKey: key, KeyID: "k", TeamID: "t", IssuedAt: time.Now().Unix() - token.TokenTimeout - 10}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			tk.GenerateIfExpired()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			_ = tk.Expired()
		}
	}()
	wg.Wait()
}
