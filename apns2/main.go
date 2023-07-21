package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ringsaturn/apns2"
	"github.com/ringsaturn/apns2/certificate"
)

var (
	certificatePath = flag.String("certificate-path", "", "Path to certificate file.")
	topic           = flag.String("topic", "", "The topic of the remote notification, which is typically the bundle ID for your app")
	mode            = flag.String("mode", "production", "APNS server to send notifications to. `production` or `development`. Defaults to `production`")
)

func main() {
	flag.Parse()

	cert, pemErr := certificate.FromPemFile(*certificatePath, "")

	if pemErr != nil {
		log.Fatalf("Error retrieving certificate `%v`: %v", certificatePath, pemErr)
	}

	client := apns2.NewClient(cert)

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
			fmt.Printf("%v: '%v'\n", res.StatusCode, res.Reason)
		}
	}
}
