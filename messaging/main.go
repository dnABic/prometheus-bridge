package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

type Options map[string]interface{}

type Stream interface {
	Connect(uri string, opts Options)
	Close()
	Publish(msg []byte, opts interface{}) error
	Consume(opts interface{}) ([]byte, error)
}

type RabbitMQStream struct {
	connection *amqp.Connection
}

type RabbitMQPublishSettings struct {
	QueueName string
}

type RabbitMQConsumeSettings struct {
	QueueName string
}

func (s *RabbitMQStream) Connect(uri string, o Options) {
	close := make(chan *amqp.Error)
	go func() {
		err := <-close

		log.Println("connection intruppted: " + err.Error())
		s.Connect(uri, o)
	}()

	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Panic(err)
	}

	s.connection = conn
	conn.NotifyClose(close)
}

func (s *RabbitMQStream) Close() {
	s.connection.Close()
}

func (s *RabbitMQStream) Publish(msg []byte, opts interface{}) error {
	options, ok := opts.(*RabbitMQPublishSettings)
	if !ok {
		log.Println("Wrong RabbitMQ publish settings, use *RabbitMQPublishSettings instead!")
	}

	ch, err := s.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ensureQueue(ch, options.QueueName)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        msg,
		})

	return err
}

func (s *RabbitMQStream) Consume(opts interface{}) ([]byte, error) {
	options, ok := opts.(*RabbitMQConsumeSettings)
	if !ok {
		log.Println("Wrong RabbitMQ publish settings, use *RabbitMQPublishSettings instead!")
	}

	ch, err := s.connection.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	q, err := ensureQueue(ch, options.QueueName)
	if err != nil {
		return nil, err
	}

	delivery, ok, err := ch.Get(q.Name, true)
	if ok {
		return delivery.Body, nil
	}

	return nil, err
}

func ensureQueue(ch *amqp.Channel, q string) (*amqp.Queue, error) {
	qu, err := ch.QueueDeclare(
		q,     // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return &qu, err
}

var _ Stream = (*RabbitMQStream)(nil)
