package messaging

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Options map[string]interface{}

type Consumer interface {
	Consume(opts interface{}) ([][]byte, error)
}

type Publisher interface {
	Publish(msg []byte, opts interface{}) error
}

type Connector interface {
	Connect(uri string, opts Options)
	Close()
}

type Stream interface {
	Connector
	Publisher
	Consumer
}

type RabbitMQConnector struct {
	connection *amqp.Connection
}

type RabbitMQStream struct {
	RabbitMQConnector
}

type RabbitMQPublishSettings struct {
	QueueName string
}

type RabbitMQConsumeSettings struct {
	QueueName    string
	MessageCount uint
}

var collector = initializeMetrics()

func NewRabbitMQStream(uri string, o Options) Stream {
	s := &RabbitMQStream{RabbitMQConnector{}}
	s.Connect(uri, o)

	return s
}

func (s *RabbitMQConnector) Connect(uri string, o Options) {
	close := make(chan *amqp.Error)
	go func() {
		err := <-close

		log.Println("connection intruppted: " + err.Error())
		s.Connect(uri, o)
	}()

	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Println("Could not initiate connection to AMQP: ", err)
		collector.BrokerConnections.WithLabelValues(uri, ConnectionFailed).Inc()

		time.Sleep(time.Second)
		s.Connect(uri, o)

		return
	}

	collector.BrokerConnections.WithLabelValues(uri, ConnectionSucceeded).Inc()

	s.connection = conn
	conn.NotifyClose(close)
}

func (s *RabbitMQConnector) Close() {
	s.connection.Close()
}

func (s *RabbitMQStream) Publish(msg []byte, opts interface{}) error {
	options, ok := opts.(*RabbitMQPublishSettings)
	if !ok {
		log.Println("Wrong RabbitMQ publish settings, use *RabbitMQPublishSettings instead!")
	}

	defer measureElapsedTime(time.Now(), collector.PublishMessageSummary)

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

func (s *RabbitMQStream) Consume(opts interface{}) ([][]byte, error) {
	options, ok := opts.(*RabbitMQConsumeSettings)
	if !ok {
		log.Println("Wrong RabbitMQ publish settings, use *RabbitMQPublishSettings instead!")
	}

	defer measureElapsedTime(time.Now(), collector.ConsumeMessageSummary)

	ch, err := s.connection.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	q, err := ensureQueue(ch, options.QueueName)
	if err != nil {
		return nil, err
	}

	if options.MessageCount < 1 {
		options.MessageCount = 1
	}
	result := make([][]byte, options.MessageCount)

	for i := uint(0); i < options.MessageCount; i++ {
		delivery, ok, err := ch.Get(q.Name, true)
		if !ok {
			return result, err
		}

		result[i] = delivery.Body
	}

	return result, nil
}

func ensureQueue(ch *amqp.Channel, q string) (*amqp.Queue, error) {
	qu, err := ch.QueueDeclare(
		q,     // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return &qu, err
}

var _ Stream = (*RabbitMQStream)(nil)
