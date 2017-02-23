package certificate

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PKCS#12

func TestParsePrivateKey(t *testing.T) {
	mockErr := errors.New("Fake-Error")
	success1Key := &rsa.PrivateKey{}
	success8Key := &rsa.PrivateKey{}
	fail1Key := &rsa.PrivateKey{}
	fail8Key := &rsa.PrivateKey{}

	// mock out private key parsing. No need to test that x509.ParsePKCS* funcs work
	// These mocks will still work with real certificate data to call x509 parse functions
	ParsePKCS8PrivateKey = func(der []byte) (interface{}, error) {
		if string(der) == "testpkcs8" {
			return success8Key, nil
		}
		if string(der)[:9] == "test-fail" {
			return fail8Key, mockErr
		}
		return x509.ParsePKCS8PrivateKey(der)
	}

	ParsePKCS1PrivateKey = func(der []byte) (*rsa.PrivateKey, error) {
		if string(der) == "testpkcs1" {
			return success1Key, nil
		}
		if string(der)[:4] == "test" {
			return fail1Key, mockErr
		}
		return x509.ParsePKCS1PrivateKey(der)
	}

	defer func() {
		ParsePKCS1PrivateKey = x509.ParsePKCS1PrivateKey
		ParsePKCS8PrivateKey = x509.ParsePKCS8PrivateKey
	}()

	type Test struct {
		input  string
		output interface{}
		err    error
	}

	tests := []Test{
		{input: "testpkcs1", output: success1Key, err: nil},
		{input: "testpkcs8", output: success8Key, err: nil},
		{input: "test-failpkcs1", output: nil, err: ErrFailedToParsePKCS1PrivateKey},
		{input: "test-failpkcs8", output: nil, err: ErrFailedToParsePKCS1PrivateKey},
	}

	for _, test := range tests {
		key, err := parsePrivateKey([]byte(test.input))
		if err != test.err {
			t.Fatal("Unexpected error", err, test.err)
		}
		if key != test.output {
			t.Fatal("Unexpected key", key, test.output)
		}
	}

}

func TestValidCertificateFromP12File(t *testing.T) {
	cer, err := FromP12File("_fixtures/certificate-valid.p12", "")
	assert.Nil(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestValidCertificateFromP12Bytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.p12")
	cer, err := FromP12Bytes(bytes, "")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestEncryptedValidCertificateFromP12File(t *testing.T) {
	cer, err := FromP12File("_fixtures/certificate-valid-encrypted.p12", "password")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestNoSuchFileP12File(t *testing.T) {
	cer, err := FromP12File("", "")
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
	assert.Equal(t, tls.Certificate{}, cer)
}

func TestBadPasswordP12File(t *testing.T) {
	cer, err := FromP12File("_fixtures/certificate-valid-encrypted.p12", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, errors.New("pkcs12: decryption password incorrect").Error(), err.Error())
}

// PEM

func TestValidCertificateFromPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid.pem", "")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestValidCertificateFromPemBytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.pem")
	cer, err := FromPemBytes(bytes, "")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestEncryptedValidCertificateFromPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid-encrypted.pem", "password")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestNoSuchFilePemFile(t *testing.T) {
	cer, err := FromPemFile("", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
}

func TestBadPasswordPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-valid-encrypted.pem", "badpassword")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrFailedToDecryptKey, err)
}

func TestBadKeyPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-bad-key.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrFailedToParsePKCS1PrivateKey, err)
}

func TestNoKeyPemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-no-key.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrNoPrivateKey, err)
}

func TestNoCertificatePemFile(t *testing.T) {
	cer, err := FromPemFile("_fixtures/certificate-no-pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, ErrNoCertificate, err)
}

func verifyHostname(cert tls.Certificate) error {
	if cert.Leaf == nil {
		return errors.New("expected leaf cert")
	}
	return cert.Leaf.VerifyHostname("APNS/2 Development IOS Push Services: com.sideshow.Apns2")
}
