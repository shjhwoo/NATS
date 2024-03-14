package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	url := "my-nats:4222"
	nc, _ := nats.Connect(url)

	defer nc.Drain()

	sub, _ := nc.Subscribe("greet.*", func(msg *nats.Msg) {

		name := msg.Subject[6:]
		msg.Respond([]byte("hello, " + name))
	})

	rep, _ := nc.Request("greet.joe", nil, time.Second)
	fmt.Println(string(rep.Data))

	rep, _ = nc.Request("greet.sue", nil, time.Second)
	fmt.Println(string(rep.Data))

	rep, _ = nc.Request("greet.bob", nil, time.Second)

	fmt.Println(string(rep.Data))

	sub.Unsubscribe()

	_, err := nc.Request("greet.joe", nil, time.Second)
	fmt.Println(err)

	//종료 방지를 위한 코드
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
}
