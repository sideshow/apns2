package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"gopkg.in/alecthomas/kingpin.v2"
)

const concurrentModeHelpText = `
Use this flag to send requests concurrently.
You should pass a Message ID to identify each request.

Expected input:  <message-id> <token> <payload>
Expected output: <message-id> <response code> <reason>

* Note that for bad requests APNs won't return an <apns-id>, that's why we need an arbitrary message-id to be passed.

Examples:
	IN: 1 aff0c63d9eaa63ad161bafee732d5bc2c31f66d552054718ff19ce314371e5d0 {"aps": {"alert": "hi"}}
	OUT: 1 200
	IN: 2 aff0c63d9eaa63ad161bafee732d5bc2c31f66d552054718ff19ce314371e5d012 {"aps": {"alert": "hi"}}
	OUT: 2 400 BadToken

Use --worker-pool-size to control the number of goroutines to be used.
`

var (
	certificatePath = kingpin.Flag("certificate-path", "Path to certificate file.").Required().Short('c').String()
	topic           = kingpin.Flag("topic", "The topic of the remote notification, which is typically the bundle ID for your app").Required().Short('t').String()
	mode            = kingpin.Flag("mode", "APNS server to send notifications to. `production` or `development`. Defaults to `production`").Default("production").Short('m').String()
	concurrentMode  = kingpin.Flag("concurrent-mode", concurrentModeHelpText).Bool()
	workerPoolSize  = kingpin.Flag("worker-pool-size", "Max number of simultaneous requests (number of goroutines).").Default("10").Int()
)

type PushRequest struct {
	messageID    string
	notification *apns2.Notification
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Alisson Sales")
	kingpin.CommandLine.Help = `Listens to STDIN to send nofitications and writes APNS response code and reason to STDOUT.
	The expected input format is: <DeviceToken> <APNS Payload>.
	The output format: <response code> <reason>.

	Examples:
		IN: aff0c63d9eaa63ad161bafee732d5bc2c31f66d552054718ff19ce314371e5d0 {"aps": {"alert": "hi"}}
		OUT: 200
		IN: abc {"aps": {"alert": "hi"}}
		OUT: 400 BadDeviceToken`
	kingpin.Parse()

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

	var wg *sync.WaitGroup
	var pushRequests chan *PushRequest

	if *concurrentMode {
		pushRequests = make(chan *PushRequest)

		for i := 0; i < *workerPoolSize; i++ {
			go pushWorker(client, pushRequests, wg)
		}
	}

	for scanner.Scan() {
		in := scanner.Text()

		if *concurrentMode {
			notificationArgs := strings.SplitN(in, " ", 3)
			messageID := notificationArgs[0]
			token := notificationArgs[1]
			payload := notificationArgs[2]

			notification := &apns2.Notification{
				Topic:       *topic,
				DeviceToken: token,
				Payload:     payload,
			}

			pushRequests <- &PushRequest{messageID, notification}
		} else {
			notificationArgs := strings.SplitN(in, " ", 2)
			token := notificationArgs[0]
			payload := notificationArgs[1]

			res, err := sendNotification(client, *topic, token, payload, "")

			if err != nil {
				log.Fatal("Error: ", err)
			} else {
				printResponse(res)
			}
		}
	}
}

func pushWorker(client *apns2.Client, pushes <-chan *PushRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	for pushRequest := range pushes {
		res, err := client.Push(pushRequest.notification)

		if err != nil {
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("%v %v %v\n", pushRequest.messageID, res.StatusCode, res.Reason)
		}
	}
}

func sendNotification(client *apns2.Client, topic, token, payload, apnsID string) (*apns2.Response, error) {
	notification := &apns2.Notification{
		Topic:       topic,
		DeviceToken: token,
		Payload:     payload,
	}

	if apnsID != "" {
		notification.ApnsID = apnsID
	}

	return client.Push(notification)
}

func printResponse(res *apns2.Response) {
	fmt.Printf("%v %v\n", res.StatusCode, res.Reason)
}
