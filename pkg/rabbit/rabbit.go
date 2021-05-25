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
	queueName := fmt.Sprintf("queue.call.%s", b.bindingKey)
	queue, err := b.channel.QueueDeclare(
		queueName, false, true, false, false, nil)
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

	err = b.channel.QueueBind(queue.Name, b.bindingKey, exchangeName, false, nil)
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
		ReplyTo:         b.bindingKey,
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

	consumerName := fmt.Sprintf("consumer.call.%s", b.bindingKey)
	msgChan, err := b.channel.Consume(
		queue.Name, consumerName, true, false, false, false, nil)
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
	queueName := fmt.Sprintf("queue.subscribe.%s", b.bindingKey)
	queue, err := b.channel.QueueDeclare(
		queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = b.channel.QueueBind(queue.Name, b.bindingKey, exchangeName, false, nil)
	if err != nil {
		return err
	}

	consumerName := fmt.Sprintf("consumer.subscribe.%s", b.bindingKey)
	msgChan, err := b.channel.Consume(
		queue.Name, consumerName, false, false, false, false, nil)
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
