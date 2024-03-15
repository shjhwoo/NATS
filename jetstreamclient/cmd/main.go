package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/nats-io/nats.go"

	"github.com/gin-gonic/gin"
)

var Js nats.JetStreamContext

func main() {
	fmt.Println("starting jetStream client...")

	// connect to NATS (https://cloud.synadia.com/ 에서 모니터링 가능)
	// nc, _ := nats.Connect("connect.ngs.global", nats.UserCredentials("./NGS-Default-CLI.creds"), nats.Name("Test Chat App"))

	nc, err := nats.Connect("my-nats:4222")
	if err != nil {
		fmt.Println("error connecting to nats", err)
		return
	}

	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		fmt.Println("error getting jetstream context", err)
		return
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "MSG",
		Subjects: []string{"msg.>"},
	})

	js.Publish("msg.testJsInit", []byte(`{"text":"OK"}`))

	Js = js

	InitRouter()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
}

// func printStreamState(ctx context.Context, stream jetstream.Stream) {
// 	info, err := stream.Info(ctx)
// 	if err != nil {
// 		fmt.Println("error getting stream info", err)
// 		return
// 	}
// 	b, err := json.MarshalIndent(info.State, "", " ")
// 	if err != nil {
// 		fmt.Println("error marshalling stream info", err)
// 		return
// 	}
// 	fmt.Println("inspecting stream info")
// 	fmt.Println(string(b))
// }

func InitRouter() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
		AllowHeaders:     []string{"withCredentials", "Content-Type"},
		MaxAge:           0,
	}))

	r.GET("/messages/history", GetAllMessages)

	r.Run(":9090")
}

func GetAllMessages(c *gin.Context) {
	consumerName := "handler-2"
	Js.AddConsumer("MSG", &nats.ConsumerConfig{
		Durable:        consumerName,
		DeliverSubject: "handler-2",
		AckPolicy:      nats.AckExplicitPolicy,
		AckWait:        time.Second,
	})

	sub, _ := Js.SubscribeSync("msg.>", nats.Bind("MSG", consumerName))

	var messages []string
	for {
		msg, err := sub.NextMsg(time.Second)
		if err != nil {
			fmt.Println("error getting message: ", err)
			break
		}

		fmt.Printf("received %q\n, %s", msg.Subject, string(msg.Data))
		messages = append(messages, string(msg.Data))
	}

	/*
		Once you have iterated over all the messages in the stream with the consumer,
		you can get them again by simply creating a new consumer
		or by deleting that consumer (nats consumer rm) and re-creating it (nats consumer add).
	*/
	defer sub.Unsubscribe() //지우고 매번 다시 만들어야, 메세지 내용을 다시 불러올 수 있음 주의하자 컨슈머는 일회성!!
	c.JSON(200, messages)
}
