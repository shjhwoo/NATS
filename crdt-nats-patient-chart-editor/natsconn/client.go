package natsconn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var NATSClient *nats.Conn
var JetStreamClient jetstream.JetStream
var ChartStream jetstream.Stream

func PrepareChartCRDTEditorStream(url string) error {
	err := ConnectNATS(url)
	if err != nil {
		return err
	}

	err = CreateJetStreamClient()
	if err != nil {
		return err
	}

	err = CreateStream()
	if err != nil {
		return err
	}

	return nil
}

func ConnectNATS(url string) error {
	nc, err := nats.Connect(url)
	if err != nil {
		return err
	}

	fmt.Println("Connected to NATS!")
	NATSClient = nc
	return nil
}

// consultation에 대한 jetstream을 생성해야 한다
// 요양기관 수준에서 분리. 멀티테넌시를 어떻게 관리?
// 실제로는 management DB 긁어와서 해야하지만 일단 지금은 데모니까 s00001 하나만 있다고 가정하고 구현
func CreateJetStreamClient() error {
	js, err := jetstream.New(NATSClient)
	if err != nil {
		fmt.Println("error getting jetstream context", err)
		return err
	}

	fmt.Println("Connected to JetStream!")
	JetStreamClient = js
	return nil
}

func CreateStream() error {
	jetCtx := context.Background()
	stream, err := JetStreamClient.CreateStream(jetCtx, jetstream.StreamConfig{
		Name:     "CHART",
		Subjects: []string{"consultation.*"},
	})

	if err != nil {
		fmt.Println("error creating stream", err)
		return err
	}

	fmt.Println("Created Chart Stream!")
	ChartStream = stream
	return nil
}

func CreateChartEventSubscriber() {
	chartEvSubscriber, err := NATSClient.Subscribe("chart", func(message *nats.Msg) {
		fmt.Println("received message:", (message.Data))

		editedChart, err := CRDTEditChart(message.Data)
		if err != nil {
			fmt.Println("error editing chart:", err)
			return
		}

		editedChartBytes, err := json.Marshal(editedChart)
		if err != nil {
			fmt.Println("error marshalling edited chart:", err)
			return
		}

		message.Respond(editedChartBytes)
	})
	if err != nil {
		fmt.Println("error subscribing to chart.*", err)
		return
	}
	fmt.Println(chartEvSubscriber, "chart SUB")
}
