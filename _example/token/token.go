package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

func main() {
	authKeyPath := flag.String("authKey", "", "Path to .p8 APNSAuthKey file (Required)")
	deviceToken := flag.String("token", "", "Push token (Required)")
	topic := flag.String("topic", "", "Topic (Required)")
	keyID := flag.String("keyID", "", "APNS KeyID (Required)")
	teamID := flag.String("teamID", "", "APNS TeamID (Required)")
	flag.Parse()

	if *authKeyPath == "" || *deviceToken == "" || *topic == "" || *keyID == "" || *teamID == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	authKey, err := token.AuthKeyFromFile(*authKeyPath)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		KeyID:   *keyID,
		TeamID:  *teamID,
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = *deviceToken
	notification.Topic = *topic
	notification.Payload = []byte(`{
			"aps" : {
				"alert" : "Hello!"
			}
		}
	`)

	client := apns2.NewTokenClient(token).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
