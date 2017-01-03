package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

func main() {
	certPath := flag.String("cert", "", "Path to .p12 certificate file (Required)")
	count := flag.Int("count", 200, "Number of pushes to send")
	token := flag.String("token", "", "Push token (Required)")
	topic := flag.String("topic", "", "Topic (Required)")
	flag.Parse()

	if *certPath == "" || *token == "" || *topic == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cert, err := certificate.FromP12File(*certPath, "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notifications := make(chan *apns2.Notification, 100)
	responses := make(chan *apns2.Response, *count)

	client := apns2.NewClient(cert).Production()

	for i := 0; i < 50; i++ {
		go worker(client, notifications, responses)
	}

	for i := 0; i < *count; i++ {
		n := &apns2.Notification{
			DeviceToken: *token,
			Topic:       *topic,
			Payload:     payload.NewPayload().Alert(fmt.Sprintf("Hello! %v", i)),
		}
		notifications <- n
	}

	for i := 0; i < *count; i++ {
		res := <-responses
		fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}

	close(notifications)
	close(responses)
}

func worker(client *apns2.Client, notifications <-chan *apns2.Notification, responses chan<- *apns2.Response) {
	for n := range notifications {
		res, err := client.Push(n)
		if err != nil {
			log.Fatal("Push Error:", err)
		}
		responses <- res
	}
}
