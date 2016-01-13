# APNS/2

APNS/2 is an (Experimental) Golang package designed for simple, flexible and fast Apple Push Notifications on iOS, OSX and Safari using the new HTTP/2 Push provider API.

[![Build Status](https://travis-ci.org/sideshow/apns2.svg?branch=master)](https://travis-ci.org/sideshow/apns2)

## Features

- Uses new Apple APNS HTTP/2 connection
- Supports persistent connections to APNS
- Fast, modular & easy to use

## Install

1. Make sure you are running version `go1.6beta1` or later from [here](https://golang.org/dl/)
2. `go get -u golang.org/x/crypto/pkcs12`

## Example

```go
package main

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"log"
)

func main() {

	cert, pemErr := certificate.FromPemFile("../cert.pem", "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns.Notification{}
	notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	notification.Topic = "com.sideshow.Apns2"
	notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	
	client := apns.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("APNS Sent:", res.ApnsId)
}

```