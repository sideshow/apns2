// Package certificate is a helper package to make working with Apple APNS
// certificates easier - It contains functions to help load Apple APNS `.p12`
// or `.pem` files from either an in memory byte array or from a local file.
//
// To use this package, you should first get the correct Apple APNS SSL
// certificates from the Apple Developer Member Center.
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

// Possible errors when parsing a certificate.
var (
	ErrFailedToParseCert            = errors.New("failed to parse certificate PEM data")
	ErrFailedToDecryptKey           = errors.New("failed to decrypt Private Key")
	ErrFailedToParsePKCS1PrivateKey = errors.New("failed to parse PKCS1 Private Key")
)

// FromPemFile loads a `.pem` certificate from a local file and returns a
// tls.Certificate. This function is similar to the crypto/tls LoadX509KeyPair
// function, however it supports `.pem` files with the cert and key combined
// in the same file, as well as password protected key files which are both
// common with APNS certificates.
//
// Use "" as the password argument if the pem certificate is not password
// protected.
func FromPemFile(filename string, password string) (tls.Certificate, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return decodePem(bytes, password)
}

// FromPemFile loads a `.pem` certificate from an in memory byte array and
// returns a tls.Certificate. This function is similar to the crypto/tls
// X509KeyPair function, however it supports `.pem` files with the cert and
// key combined, as well as password protected keys which are both common with
// APNS certificates.
//
// Use "" as the password argument if the pem certificate is not password
// protected.
func FromPemBytes(bytes []byte, password string) (tls.Certificate, error) {
	return decodePem(bytes, password)
}

// FromP12File loads a `.p12` certificate from a local file and returns a
// tls.Certificate.
//
// Use "" as the password argument if the pem certificate is not password
// protected.
func FromP12File(filename string, password string) (tls.Certificate, error) {
	p12bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}
	return FromP12bytes(p12bytes, password)
}

// FromP12File loads a `.p12` certificate from an in memory byte array and
// returns a tls.Certificate.
//
// Use "" as the password argument if the pem certificate is not password
// protected.
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
