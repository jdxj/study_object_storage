package rabbit

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var (
	b1 *Broker
	b2 *Broker

	bk1 = "test_binding_key1"
	bk2 = "test_binding_key2"
)

func TestMain(m *testing.M) {

	b1 = New("guest", "guest", "127.0.0.1", bk1, 5672)
	err := b1.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer b1.Close()

	b2 = New("guest", "guest", "127.0.0.1", bk2, 5672)
	err = b2.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer b2.Close()

	h := func(msg *Message) (*Message, error) {
		msg.Body = append(msg.Body, []byte(", world!")...)
		return msg, nil
	}
	err = b2.Subscribe(h)
	if err != nil {
		log.Fatalln()
	}

	code := m.Run()
	os.Exit(code)
}

func TestBroker_Call(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &Message{Body: []byte("hello")}
	result, err := b1.Call(ctx, bk2, msg)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	for _, msg := range result {
		fmt.Printf("-----%s\n", msg.Body)
	}
}
