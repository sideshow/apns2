package certificate_test

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/sideshow/apns2/certificate"
	"github.com/stretchr/testify/assert"
)

// PKCS#12

func TestValidCertificateFromP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid.p12", "")
	assert.Nil(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestValidCertificateFromP12Bytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.p12")
	cer, err := certificate.FromP12Bytes(bytes, "")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestEncryptedValidCertificateFromP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid-encrypted.p12", "password")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestNoSuchFileP12File(t *testing.T) {
	cer, err := certificate.FromP12File("", "")
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
	assert.Equal(t, tls.Certificate{}, cer)
}

func TestBadPasswordP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid-encrypted.p12", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, errors.New("pkcs12: decryption password incorrect").Error(), err.Error())
}

// PEM

func TestValidCertificateFromPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid.pem", "")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestValidCertificateFromPemBytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.pem")
	cer, err := certificate.FromPemBytes(bytes, "")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestEncryptedValidCertificateFromPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-encrypted.pem", "password")
	assert.NoError(t, err)
	assert.Nil(t, verifyHostname(cer))
}

func TestNoSuchFilePemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, errors.New("open : no such file or directory").Error(), err.Error())
}

func TestBadPasswordPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-encrypted.pem", "badpassword")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, certificate.ErrFailedToDecryptKey, err)
}

func TestBadKeyPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-bad-key.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, certificate.ErrFailedToParsePKCS1PrivateKey, err)
}

func TestNoKeyPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-no-key.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, certificate.ErrNoPrivateKey, err)
}

func TestNoCertificatePemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-no-certificate.pem", "")
	assert.Equal(t, tls.Certificate{}, cer)
	assert.Equal(t, certificate.ErrNoCertificate, err)
}

func verifyHostname(cert tls.Certificate) error {
	if cert.Leaf == nil {
		return errors.New("expected leaf cert")
	}
	return cert.Leaf.VerifyHostname("APNS/2 Development IOS Push Services: com.sideshow.Apns2")
}
