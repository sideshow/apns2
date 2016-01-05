package main

import (
	apns "github.com/sideshow/apns2"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {

	cert, pemErr := apns.FromPemFile("../cert.pem", "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns.Notification{}
	notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	notification.Topic = "com.sideshow.Apns2"
	notification.Payload = []byte(`{
		  "aps" : {
			"alert" : "Hello!"
		  }
		}
	`)

	client := apns.NewClient(cert).Development()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("APNS Error: ", err)
		return
	}

	log.Println("APNS Sent: ", res.NotificationID)
}
