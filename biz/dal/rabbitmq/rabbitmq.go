package rabbitmq

import (
	"context"
	"ququiz/lintang/scoring-service/config"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

const quizQueryConsumerName = "quizQueryConsumer"

func NewRabbitMQ(cfg *config.Config) *RabbitMQ {
	zap.L().Info("rmq address: " + cfg.RabbitMQ.RMQAddress)

	conn, err := amqp.Dial(cfg.RabbitMQ.RMQAddress)

	if err != nil {
		zap.L().Fatal("error: cannot connect to rabbitmq: " + err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("error can't get rabbitmq cahnnel: " + err.Error())
	}

	err = channel.ExchangeDeclare(
		"scoring-notification",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("err: channel.ExchangeDeclare : scoring-notification " + err.Error())
	}

	err = channel.Qos(
		1, 0,
		false,
	)
	if err != nil {
		zap.L().Fatal("err: channel.Qos" + err.Error())
	}

	channel.QueueDeclare(
		"user-add-score",
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	channel.QueueBind(
		"user-add-score",
		"correct-answer",
		"scoring-quiz-query",
		false,
		nil,
	)
	// error duplicate queue gak usah dihandle
	// if err != nil {
	// 	zap.L().Fatal("cant create new queue user-add-score (r.rmq.Channel.QueueDeclare) (ListenAndServe) (RMQConsumer) ", zap.Error(err))

	// }

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

}

func (r *RabbitMQ) Close(ctx context.Context) {

	r.Connection.Close()
}
