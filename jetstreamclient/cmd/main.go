package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/gin-gonic/gin"
)

var JsConsumer jetstream.Consumer

func main() {
	fmt.Println("starting jetStream client...")

	// connect to NATS (https://cloud.synadia.com/ 에서 모니터링 가능)
	// nc, _ := nats.Connect("connect.ngs.global", nats.UserCredentials("./NGS-Default-CLI.creds"), nats.Name("Test Chat App"))

	nc, _ := nats.Connect("my-nats:4222")

	defer nc.Drain()

	js, _ := jetstream.New(nc)

	cfg := jetstream.StreamConfig{
		Name:     "MSG",
		Subjects: []string{"msg.>"},
		Storage:  jetstream.FileStorage,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, _ := js.CreateStream(ctx, cfg)
	fmt.Println("created the stream: ", stream)

	js.Publish(ctx, "msg.testJsInit", []byte(`{"text":"OK"}`))
	printStreamState(ctx, stream)

	//테스트 용도로..
	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{})
	if err != nil {
		fmt.Println("error creating consumer", err)
	}

	JsConsumer = cons

	InitRouter()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
}

func printStreamState(ctx context.Context, stream jetstream.Stream) {
	info, err := stream.Info(ctx)
	if err != nil {
		fmt.Println("error getting stream info", err)
		return
	}
	b, err := json.MarshalIndent(info.State, "", " ")
	if err != nil {
		fmt.Println("error marshalling stream info", err)
		return
	}
	fmt.Println("inspecting stream info")
	fmt.Println(string(b))
}

func InitRouter() {
	r := gin.Default()

	r.GET("/messages/history", GetAllMessages)

	r.Run(":9090")
}

func GetAllMessages(c *gin.Context) {
	JsConsumer.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		fmt.Println("received msg on", msg.Subject(), "cotent: ", string(msg.Data()))
	})
}
