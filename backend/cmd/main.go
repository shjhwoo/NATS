package main

import (
	"backend/natClient"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	NATS_URL string `default:"my-nats:4222"`
}

type App struct {
	Ctx        context.Context
	Env        Environment
	CancelFunc func()
}

var app *App
var lock = &sync.Mutex{}

func main() {
	ctx, cf := context.WithCancel(context.Background())

	app = GetInstance(ctx, cf)

	endConsume := make(chan bool)

	go app.ShutdownService()
	go app.StartService(endConsume)

	//아래 코드가 실행되는 순간 서버는 종료된다.
	<-endConsume

	<-app.Ctx.Done()
}

func GetInstance(ctx context.Context, cf context.CancelFunc) *App {
	if app == nil {
		lock.Lock()
		defer lock.Unlock()
		if app == nil {

			var env Environment
			err := envconfig.Process("", &env)
			if err != nil {
				fmt.Println("Error processing environment variables")
				panic(err)
			}

			app = &App{
				Env:        env,
				Ctx:        ctx,
				CancelFunc: cf,
			}

			natClient.ConnectNATS(env.NATS_URL)

			fmt.Println("Connected to NATS at :", env.NATS_URL)
		}
	}

	fmt.Println("App instance created", app)

	return app
}

func (app *App) StartService(endConsume chan bool) {
	// Start the service
	close(endConsume)
}

func (app *App) ShutdownService() {
	// Shutdown the service
	sigterm := make(chan os.Signal, 1)
	// syscall.SIGINT, syscall.SIGTERM 인 경우 공지한다.
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-app.Ctx.Done():
		fmt.Println("Terminating with context done")
		app.endService()
	case <-sigterm:
		fmt.Println("Terminating via system signal", nil)
		app.endService()
		app.CancelFunc() //contxt 종료시킨다.

		os.Exit(0)
	}
}

func (app *App) endService() {
	//app.Router.Shutdown(app.Ctx)
}
