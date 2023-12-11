package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lgaches/apns2"
	"github.com/lgaches/apns2/certificate"
	"github.com/lgaches/apns2/token"
)

func main() {
	certificatePath := flag.String("certificate-path", "", "Path to certificate file.")
	authKeyPath := flag.String("authkey-path", "", "path to the P8 file. (Certificates, Identifiers & Profiles -> Keys)")
	keyID := flag.String("key-id", "", "Key ID from developer account (Certificates, Identifiers & Profiles -> Keys)")
	teamID := flag.String("team-id", "", "Team ID from developer account (View Account -> Membership)")
	topic := flag.String("topic", "", "The topic of the remote notification, which is typically the bundle ID for your app")
	mode := flag.String("mode", "production", "APNS server to send notifications to. `production` or `development`. Defaults to `production`")

	flag.Parse()

	var client *apns2.Client

	if certificatePath == nil || *certificatePath != "" {
		cert, pemErr := certificate.FromPemFile(*certificatePath, "")
		if pemErr != nil {
			log.Fatalf("Error retrieving certificate `%v`: %v", certificatePath, pemErr)
		}
		client = apns2.NewClient(cert)
	} else if *authKeyPath != "" || *teamID != "" || *keyID != "" {
		authKey, authErr := token.AuthKeyFromFile(*authKeyPath)

		if authErr != nil {
			log.Fatalf("Error retrieving Token `%v`: %v", authKeyPath, authErr)
		}

		authToken := &token.Token{
			AuthKey: authKey,
			KeyID:   *keyID,
			TeamID:  *teamID,
		}

		client = apns2.NewTokenClient(authToken)
	} else {
		flag.Usage()
		os.Exit(1)
	}

	if *topic == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *mode == "development" {
		client.Development()
	} else {
		client.Production()
	}

	fmt.Println("Ready to send push notifications. Enter tokens and payloads then press enter. Example: aff0c63d9eaa63ad161bafee732d5bc2c31f66d552054718ff19ce314371e5d0 {\"aps\": {\"alert\": \"hi\"}}")
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
