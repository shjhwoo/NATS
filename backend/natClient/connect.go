package natClient

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var natsClient *nats.Conn

var connectCount int

func ConnectNATS(natsURL string) {
	if connectCount >= 5 {
		panic("Max retries reached. Exiting...")
	}

	if natsClient != nil {
		fmt.Println("Connected to NATS at :", natsURL)
		return
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		fmt.Println("Error connecting to NATS at :", natsURL)
		connectCount++
		fmt.Println("retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		ConnectNATS(natsURL)
	}

	natsClient = nc
}
