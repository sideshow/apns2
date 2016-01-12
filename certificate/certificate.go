package certificate

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"strings"
)

var (
	ErrFailedToParseCert          = errors.New("failed to parse certificate PEM data")
	ErrFailedToDecryptKey         = errors.New("failed to decrypt Private Key")
	ErrFailedToParsePKCS1PrivateKey = errors.New("failed to parse PKCS1 Private Key")
)

func FromPemFile(filename string, password string) (tls.Certificate, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return decodePem(bytes, password)
}

func FromPemBytes(bytes []byte, password string) (tls.Certificate, error) {
	return decodePem(bytes, password)
}

func FromP12File(filename string, password string) (tls.Certificate, error) {
	p12bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return FromP12bytes(p12bytes, password)
}

func FromP12bytes(bytes []byte, password string) (tls.Certificate, error) {
	key, cert, err := pkcs12.Decode(bytes, password)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}, nil
}

func decodePem(bytes []byte, password string) (tls.Certificate, error) {
	var cert tls.Certificate
	var block *pem.Block
	for {
		block, bytes = pem.Decode(bytes)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		}
		if block.Type == "PRIVATE KEY" || strings.HasSuffix(block.Type, "PRIVATE KEY") {
			key, err := decodeKey(block, password)
			if err != nil {
				return cert, err
				break
			}
			cert.PrivateKey = key
		}
	}
	if len(cert.Certificate) == 0 {
		return cert, ErrFailedToParseCert
	}
	if c, e := x509.ParseCertificate(cert.Certificate[0]); e == nil {
		cert.Leaf = c
	}
	return cert, nil
}

func decodeKey(block *pem.Block, password string) (crypto.PrivateKey, error) {
	if x509.IsEncryptedPEMBlock(block) {
		bytes, decryptErr := x509.DecryptPEMBlock(block, []byte(password))
		if decryptErr != nil {
			return nil, ErrFailedToDecryptKey
		}
		key, parseErr := x509.ParsePKCS1PrivateKey(bytes)
		if parseErr != nil {
			return nil, ErrFailedToParsePKCS1PrivateKey
		}
		return key, nil
	}
	return block, nil
}
