package rabbitmq

import (
	"context"
	"encoding/json"
	"ququiz/lintang/scoring-service/biz/domain"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type ScoringMQ struct {
	ch *amqp.Channel
}

func NewScoringMQ(rmq *RabbitMQ) *ScoringMQ {
	return &ScoringMQ{
		ch: rmq.Channel,
	}
}

func (s *ScoringMQ) SendQuizRecap(ctx context.Context, quizRecapMessage domain.QuizRecapMessage) error  {
	return s.publish(ctx, "quiz-score-notification", quizRecapMessage)
}

func (s *ScoringMQ) publish(ctx context.Context, routingKey string, event interface{}) error {

	jsonBody, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("json.Marshal (publish) (ScoringMQ)", zap.Error(err))
		return err
	}
	zap.L().Info("send json serialized ata to notification service!!!")
	err = s.ch.Publish(
		"scoring-notification",
		routingKey, // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "scoring-service",
			ContentType: "application/json",
			Body:        jsonBody,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		zap.L().Error("m.ch.Publish: (SendQuizRecap)", zap.Error(err))
		return err
	}
	return nil
}
