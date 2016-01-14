# APNS/2
APNS/2 is an (Experimental) go package designed for simple, flexible and fast Apple Push Notifications on iOS, OSX and Safari using the new HTTP/2 Push provider API.

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
}
```

## Notification
At a minimum, a _Notification_ needs a _Token_, a _Topic_ and a _Payload_.

```go
notification := &Notification{
  Token: "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7",
  Topic: "com.sideshow.Apns2"
  Payload: []byte(`{"aps":{"alert":"Hello!"}}`),
}
```

You can also set an optional _ApnsId_, _Expiration_ or _Priority_.

```go
notification.ApnsId =  "40636A2C-C093-493E-936A-2A4333C06DEA"
notification.Expiration = time.Now()
notification.Priority = apns.PriorityLow
```

## Response, Error handling
APNS/2 draws the distinction between a valid response from Apple indicating wether or not the _Notification_ was sent or not, and an unrecoverable or unexpected _Error_;
- An `Error` is returned if a non-recoverable error occurs, i.e. if there is a problem with the underlying _http.Client_ connection or _Certificate_, the payload was not sent, or a valid _Response_ was not received.
- A `Response` is returned if the payload was successfully sent to Apple and a documented response was received. This struct will contain more information about whether or not the push notification succeeded, its _apns-id_ and if applicable, more information around why it did not succeed.

To check if a `Notification` was successfully sent;

```go
res, err := client.Push(notification)
if err != nil {
  log.Println("There was an error", err)
  return
}
if res.Sent() {
  log.Println("Successfully sent", res.ApnsId())
}
```

## License
The MIT License (MIT)

Copyright (c) 2016 Adam Jones

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
