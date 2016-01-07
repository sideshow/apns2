package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
)

func FromPemFile(filename string, password string) (tls.Certificate, error) {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return decodePem(tls.Certificate{}, pemBytes, password)
}

func FromPemBytes(bytes []byte, password string) (tls.Certificate, error) {
	return decodePem(tls.Certificate{}, bytes, password)
}

func FromP12File(filename string, password string) (tls.Certificate, error) {
	p12bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return decodeP12(p12bytes, password)
}

func decodePem(cert tls.Certificate, bytes []byte, password string) (tls.Certificate, error) {
	block, rest := pem.Decode(bytes)
	if block == nil {
		return cert, nil
	}
	if x509.IsEncryptedPEMBlock(block) {
		_, err := x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return cert, errors.New("Error decrypting certificate")
		}
	}
	switch block.Type {
	case "CERTIFICATE":
		cert.Certificate = append(cert.Certificate, block.Bytes)
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return cert, errors.New("Error parsing RSA PRIVATE KEY")
		}
		cert.PrivateKey = key
	default:
		return cert, errors.New("Cert block wasn't CERTIFICATE or PRIVATE KEY")
	}
	return decodePem(cert, rest, password)
}

func decodeP12(bytes []byte, password string) (tls.Certificate, error) {
	key, cert, err := pkcs12.Decode(bytes, password)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}, nil
}
