package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"ququiz/lintang/scoring-service/biz/domain"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type QueryQueryProducer struct {
	ch *amqp.Channel
}

func NewQuizQueryProducer(rmq *RabbitMQ) *QueryQueryProducer {
	return &QueryQueryProducer{
		ch: rmq.Channel,
	}
}

func (p *QueryQueryProducer) SendDeleteCache(ctx context.Context, deleteCacheMessage domain.DeleteCacheMessage) error {
	return p.publishToQuizQuery(ctx, "delete-cache", deleteCacheMessage)
}

func (p *QueryQueryProducer) publishToQuizQuery(ctx context.Context, routingKey string, event interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		zap.L().Error("gob.NewEncoder(&b).Encode(event) (publishToQuizQuery) (SendDeleteCache)", zap.Error(err))
		return err
	}

	zap.L().Info("send json serialized ata to quiz query service!!!")
	err := p.ch.Publish(
		"scoring-quiz-query",
		routingKey, // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "scoring-service-quiz-query",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		zap.L().Error("p.ch.Publish (publishToQuizQuery) (SendDeleteCache)", zap.Error(err))
		return err
	}

	return nil
}
