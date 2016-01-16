package certificate_test

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/sideshow/apns2/certificate"
)

// p12

func TestValidCertificateFromP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid.p12", "")
	if err != nil {
		t.Fatal(err)
	}
	if e := verifyHostname(cer); e != nil {
		t.Fatal(e)
	}
}

func TestValidCertificateFromP12Bytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.p12")
	cer, err := certificate.FromP12Bytes(bytes, "")
	if err != nil {
		t.Fatal(err)
	}
	if e := verifyHostname(cer); e != nil {
		t.Fatal(e)
	}
}

func TestEncryptedValidCertificateFromP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid-encrypted.p12", "password")
	if err != nil {
		t.Fatal(err)
	}
	if e := verifyHostname(cer); e != nil {
		t.Fatal(e)
	}
}

func TestNoSuchFileP12File(t *testing.T) {
	_, err := certificate.FromP12File("", "")
	if err.Error() != errors.New("open : no such file or directory").Error() {
		t.Fatal("expected error", "open : no such file or directory")
	}
}

func TestBadPasswordP12File(t *testing.T) {
	_, err := certificate.FromP12File("_fixtures/certificate-valid-encrypted.p12", "")
	if err.Error() != errors.New("pkcs12: decryption password incorrect").Error() {
		t.Fatal("expected", "pkcs12: decryption password incorrect")
	}
}

// pem

func TestValidCertificateFromPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid.pem", "")
	if err != nil {
		t.Fatal(err)
	}
	if e := verifyHostname(cer); e != nil {
		t.Fatal(e)
	}
}

func TestValidCertificateFromPemBytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.pem")
	cer, err := certificate.FromPemBytes(bytes, "")
	if err != nil {
		t.Fatal(err)
	}
	if e := verifyHostname(cer); e != nil {
		t.Fatal(e)
	}
}

func TestEncryptedValidCertificateFromPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-encrypted.pem", "password")
	if err != nil {
		t.Fatal(err)
	}
	if e := verifyHostname(cer); e != nil {
		t.Fatal(e)
	}
}

func TestNoSuchFilePemFile(t *testing.T) {
	_, err := certificate.FromPemFile("", "")
	if err.Error() != errors.New("open : no such file or directory").Error() {
		t.Fatal("expected error", "open : no such file or directory")
	}
}

func TestBadPasswordPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-encrypted.pem", "badpassword")
	if err != certificate.ErrFailedToDecryptKey {
		t.Fatal("expected error", certificate.ErrFailedToDecryptKey, cer)
	}
}

func TestBadKeyPemFile(t *testing.T) {
	_, err := certificate.FromPemFile("_fixtures/certificate-bad-key.pem", "")
	if err != certificate.ErrFailedToParsePKCS1PrivateKey {
		t.Fatal("expected error", certificate.ErrFailedToParsePKCS1PrivateKey)
	}
}

func TestNoKeyPemFile(t *testing.T) {
	_, err := certificate.FromPemFile("_fixtures/certificate-no-key.pem", "")
	if err != certificate.ErrNoPrivateKey {
		t.Fatal("expected error", certificate.ErrNoPrivateKey)
	}
}

func TestNoCertificatePemFile(t *testing.T) {
	_, err := certificate.FromPemFile("_fixtures/certificate-no-certificate.pem", "")
	if err != certificate.ErrNoCertificate {
		t.Fatal("expected error", certificate.ErrNoCertificate)
	}
}

func verifyHostname(cert tls.Certificate) error {
	if cert.Leaf == nil {
		return errors.New("expected leaf cert")
	}
	return cert.Leaf.VerifyHostname("APNS/2 Development IOS Push Services: com.sideshow.Apns2")
}
