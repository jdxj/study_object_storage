package rabbit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "object_storage"
	exchangeKind = "topic"
)

var (
	broker *Broker
)

func Init(user, pass, host, bindingKey string, port int) error {
	broker = New(user, pass, host, bindingKey, port)
	return broker.Connect()
}

func Call(ctx context.Context, routingKey string, msg *Message) ([]*Message, error) {
	return broker.Call(ctx, routingKey, msg)
}

func Publish(routingKey string, msg *Message) error {
	return broker.Publish(routingKey, msg)
}

func Subscribe(h Handler) error {
	return broker.Subscribe(h)
}

type Handler func(*Message) (*Message, error)

type Message struct {
	Headers map[string]interface{}
	Body    []byte
}

func New(user, pass, host, bindingKey string, port int) *Broker {
	b := &Broker{
		user:       user,
		pass:       pass,
		host:       host,
		port:       port,
		bindingKey: bindingKey,
	}

	return b
}

type Broker struct {
	user string
	pass string
	host string
	port int

	bindingKey string
	conn       *amqp.Connection
	channel    *amqp.Channel
}

func (b *Broker) Connect() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d", b.user, b.pass, b.host, b.port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	b.conn = conn
	b.channel, err = conn.Channel()
	if err != nil {
		return err
	}

	return b.channel.ExchangeDeclare(
		exchangeName, exchangeKind, true, false, false, false, nil)
}

// Call 期待响应
func (b *Broker) Call(ctx context.Context, routingKey string, msg *Message) ([]*Message, error) {
	queue, err := b.channel.QueueDeclare(
		"", false, true, false, false, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		// 调用完立即删除队列
		_, err := b.channel.QueueDelete(queue.Name, false, false, true)
		if err != nil {
			log.Printf("QueueDelete: %s", err)
		}
	}()

	err = b.channel.QueueBind(queue.Name, queue.Name, exchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	pub := amqp.Publishing{
		Headers:         msg.Headers,
		ContentType:     "",
		ContentEncoding: "",
		DeliveryMode:    1,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         queue.Name,
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            msg.Body,
	}
	err = b.channel.Publish(exchangeName, routingKey, false, false, pub)
	if err != nil {
		return nil, err
	}

	msgChan, err := b.channel.Consume(
		queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	var result []*Message
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return result, nil
			}
			return nil, ctx.Err()
		case d := <-msgChan:
			msg := &Message{
				Headers: d.Headers,
				Body:    d.Body,
			}
			result = append(result, msg)
		}
	}
}

// Publish 不需要响应, 推完就拉到
func (b *Broker) Publish(routingKey string, msg *Message) error {
	pub := amqp.Publishing{
		Headers:         msg.Headers,
		ContentType:     "",
		ContentEncoding: "",
		DeliveryMode:    1,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            msg.Body,
	}
	return b.channel.Publish(
		exchangeName, routingKey, false, false, pub)
}

func (b *Broker) Subscribe(h Handler) error {
	queue, err := b.channel.QueueDeclare(
		"", false, true, false, false, nil)
	if err != nil {
		return err
	}
	err = b.channel.QueueBind(queue.Name, b.bindingKey, exchangeName, false, nil)
	if err != nil {
		return err
	}

	msgChan, err := b.channel.Consume(
		queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			d := <-msgChan
			msgR := &Message{
				Headers: d.Headers,
				Body:    d.Body,
			}
			msgA, err := h(msgR)
			if err != nil {
				log.Printf("rabbit handler: %s\n", err)
				continue
			}

			if msgA != nil { // 有响应, 那么发送反馈
				err = b.Publish(d.ReplyTo, msgA)
				if err != nil {
					log.Printf("b.Publish: %s", err)
					continue
				}
			}

			err = d.Ack(false)
			if err != nil {
				log.Printf("rabbit ack: %s\n", err)
			}
		}
	}()
	return nil
}

func (b *Broker) Close() error {
	b.channel.Close()
	return b.conn.Close()
}
