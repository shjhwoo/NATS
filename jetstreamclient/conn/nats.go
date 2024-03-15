package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var JetStreamClient jetstream.JetStream

func StartNatsServer() {
	url := "my-nats:4222"
	nc, _ := nats.Connect(url)

	defer nc.Drain()

	js, _ := jetstream.New(nc)

	cfg := jetstream.StreamConfig{
		Name:     "CHAT",
		Subjects: []string{"message.*.>"},
		Storage:  jetstream.FileStorage,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, _ := js.CreateStream(ctx, cfg)
	fmt.Println("created the stream")

	js.Publish(ctx, "message.connected", []byte("connected"))
	printStreamState(ctx, stream)

	JetStreamClient = js
}

func printStreamState(ctx context.Context, stream jetstream.Stream) {
	info, _ := stream.Info(ctx)
	b, _ := json.MarshalIndent(info.State, "", " ")
	fmt.Println("inspecting stream info")
	fmt.Println(string(b))
}
