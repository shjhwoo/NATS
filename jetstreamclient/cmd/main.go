package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/net/context"

	"github.com/gorilla/websocket"
)

var RDB *sqlx.DB
var JetStream jetstream.JetStream
var JetStreamKV jetstream.KeyValue

func main() {

	wsc, _, err := websocket.DefaultDialer.Dial("ws://my-nats:8080", nil)
	if err != nil {
		fmt.Println("error dialing websocket server", err)
		return
	}
	go func() {
		for {
			_, message, err := wsc.ReadMessage()
			if err != nil {
				fmt.Println("error reading message from websocket server", err)
				return
			}

			fmt.Println("websocket message:", string(message))
			//연결해서 클라이언트 종료됨을 감지하면
			//DB에다가 사용자 lastAccess 시간을 저장한다.

		}
	}()

	// connect to NATS (https://cloud.synadia.com/ 에서 모니터링 가능)
	// nc, _ := nats.Connect("connect.ngs.global", nats.UserCredentials("./NGS-Default-CLI.creds"), nats.Name("Test Chat App"))

	nc, err := nats.Connect("my-nats:4222")
	if err != nil {
		fmt.Println("error connecting to nats", err)
		return
	}

	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		fmt.Println("error getting jetstream context", err)
		return
	}

	JetStream = js

	jetCtx := context.Background()
	stream, _ := JetStream.CreateStream(jetCtx, jetstream.StreamConfig{
		Name:     "MSG",
		Subjects: []string{"msg.>"},
	})

	loginSub, err := nc.Subscribe("login", func(message *nats.Msg) {
		fmt.Println("received message:", (message.Data))
		//사용자 있는지 확인하고 이상없으면 로그인 시간을 KV에 저장
		user := string(message.Data)
		err = findUser(&User{Username: user})
		if err != nil {
			fmt.Println("error finding user:", err)
			return
		}

		lastAccess, err := saveUserLoginInfoToKV(&User{Username: user})
		if err != nil {
			fmt.Println("error saving user login info to KV:", err)
			return
		}

		message.Respond([]byte(lastAccess))
	})
	if err != nil {
		fmt.Println("error subscribing to login.*", err)
		return
	}
	fmt.Println(loginSub, "login SUB")

	MSG, err := nc.Subscribe("meta.ALL", func(message *nats.Msg) {

		last1hour := time.Now().Add(-time.Hour) //최근 1시간 내의 메세지만 보여주는 기능.
		cons, err := stream.CreateOrUpdateConsumer(jetCtx, jetstream.ConsumerConfig{
			DeliverPolicy: jetstream.DeliverByStartTimePolicy, //최근 1시간 내의 메세지만 보여주는 기능.
			OptStartTime:  &last1hour,                         //최근 1시간 내의 메세지만 보여주는 기능.
			FilterSubject: "msg.>",
			Durable:       "getArchivedMessages",
		})
		if err != nil {
			fmt.Println("error creating consumer", err)
			return
		}

		freshStreamInfo, err := stream.Info(jetCtx)
		if err != nil {
			fmt.Println("error getting stream info", err)
			return
		}

		messageCount := freshStreamInfo.State.Msgs
		msgs, err := cons.FetchNoWait(int(messageCount))
		if err != nil {
			fmt.Println("error fetching messages", err)
			return
		}

		var msgList []string
		for msg := range msgs.Messages() {
			msg.Ack()
			msgList = append(msgList, string(msg.Data()))
		}

		// fmt.Println(len(msgList), "개의 메세지가 쌓여있음")
		msgListBytes, err := json.Marshal(msgList)
		if err != nil {
			fmt.Println("error marshalling message list", err)
			return
		}

		err = message.Respond(msgListBytes)
		if err != nil {
			fmt.Println("error responding to request", err)
			return
		}

		err = stream.DeleteConsumer(jetCtx, "getArchivedMessages")
		if err != nil {
			fmt.Println("error deleting consumer", err)
			return
		}
		fmt.Println("deleted consumer getArchivedMessages")
	})
	if err != nil {
		fmt.Println("error subscribing to meta.ALL", err)
		return
	}
	fmt.Println(MSG, "MSG SUB")

	KV, err := JetStream.CreateKeyValue(jetCtx, jetstream.KeyValueConfig{
		Bucket: "userLoginInfo",
	})
	if err != nil {
		fmt.Println("error creating key value store", err)
		return
	}
	JetStreamKV = KV

	ConnectRDB()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
}

func ConnectRDB() {
	// rdb 연결
	db, err := sqlx.Open("mysql", "root:1234@tcp(host.docker.internal:3306)/natsuser")
	if err != nil {
		fmt.Println("error connecting to MySQL database:", err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error pinging MySQL database:", err)
		return
	}

	fmt.Println("connected to MySQL database")
	RDB = db
}

type User struct {
	Username string `json:"username"`
}

func findUser(user *User) error {
	var username string
	err := RDB.Get(&username, "SELECT username FROM natsuser.users WHERE username = ?", user.Username)
	if err != nil {
		fmt.Println("error scanning user from database:", err)
		return err
	}
	return nil
}

func saveUserLoginInfoToKV(user *User) (string, error) {
	koreaLocation := time.FixedZone("KST", 9*60*60)
	now := time.Now().UTC()
	koreaNow := now.In(koreaLocation)
	lastAccess := koreaNow.Format("2006-01-02 15:04:05")

	res, err := RDB.Exec("UPDATE natsuser.users SET lastAccess = ? WHERE username = ?", lastAccess, user.Username)
	if err != nil {
		fmt.Println("error updating user last access time:", err)
		return "", err
	}

	fmt.Println("created user login info KV pair:", res)
	return lastAccess, nil
}
