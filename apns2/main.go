package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lgaches/apns2"
	"github.com/lgaches/apns2/token"
)

var (
	tokenPath = flag.String("token-path", "", "Path to token file.")
	teamID    = flag.String("team-id", "", "The team ID")
	keyID     = flag.String("key-id", "", "The Key ID")
	topic     = flag.String("topic", "", "The topic of the remote notification, which is typically the bundle ID for your app")
	mode      = flag.String("mode", "production", "APNS server to send notifications to. `production` or `development`. Defaults to `production`")
)

func main() {
	flag.Parse()

	authKey, authErr := token.AuthKeyFromFile(*tokenPath)

	if authErr != nil {
		log.Fatalf("Error retrieving Token `%v`: %v", tokenPath, authErr)
	}

	authToken := &token.Token{
		AuthKey: authKey,
		KeyID:   *keyID,
		TeamID:  *teamID,
	}

	client := apns2.NewTokenClient(authToken)

	if *mode == "development" {
		client.Development()
	} else {
		client.Production()
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		in := scanner.Text()
		notificationArgs := strings.SplitN(in, " ", 2)
		token := notificationArgs[0]
		payload := notificationArgs[1]

		notification := &apns2.Notification{
			DeviceToken: token,
			Topic:       *topic,
			Payload:     payload,
		}

		res, err := client.Push(notification)

		if err != nil {
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("%v: '%v' . %v - %v  - %v\n", res.StatusCode, res.Reason, res.ApnsID, res.Timestamp, res.ApnsUniqueID)
		}
	}
}
