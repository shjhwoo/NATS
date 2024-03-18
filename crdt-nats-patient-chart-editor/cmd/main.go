package main

import (
	"chart-coeditor/environment"
	"chart-coeditor/natsconn"
	"os"
	"os/signal"
)

func main() {
	//NATS 연결함
	env := environment.GetEnvironment()
	err := natsconn.PrepareChartCRDTEditorStream(env.NATS_URL)
	if err != nil {
		panic(err)
	}

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt)
	<-quitChan
}
