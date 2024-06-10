package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"ququiz/lintang/scoring-service/biz/domain"

	"go.uber.org/zap"
)

type LeaderboardRedis interface {
	CalculateUserScore(ctx context.Context, weight uint64,
		userID string, quizID string, userName string) error
	GetTopLeaderBoard(ctx context.Context, quizID string) ([]domain.RedisLeaderBoard, error)
}

type QuizQueryMQConsumer struct {
	rmq              *RabbitMQ
	leaderboardRedis LeaderboardRedis
}

func NewQuizQueryMQConsumer(r *RabbitMQ, l LeaderboardRedis) *QuizQueryMQConsumer {
	return &QuizQueryMQConsumer{r, l}
}

const rabbitMQConsumerName = "scoring-svc-consumer"

func (r *QuizQueryMQConsumer) ListenAndServe() error {
	
	err := r.rmq.Channel.QueueBind(
		"user-add-score",
		"correct-answer",
		"scoring-quiz-query",
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("cant bind queue %s to exchange scoring-quiz-query", "user-add-score"))
	}
	msgs, err := r.rmq.Channel.Consume(
		"user-add-score",
		rabbitMQConsumerName,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprint("cant consume message from queue %s", "user-add-score"))
	}

	go func() {
		for msg := range msgs {
			zap.L().Info("Received message: %s" + msg.RoutingKey)
			cAnswer, err := decodeCorrectAnswer(msg.Body)
			if err != nil {
				zap.L().Error("decodeCorrectAnswer (ListenAndServe) (QuizQueryMQConsumer) ", zap.Error(err))
				return
			}
			var nack bool
			switch msg.RoutingKey {
			case "correct-answer":
				// TODO:  implement update leaderboard
				err := r.leaderboardRedis.CalculateUserScore(context.Background(), cAnswer.Weight,
					cAnswer.UserID, cAnswer.QuizID,
					cAnswer.Username)

				if err != nil {
					zap.L().Error("r.leaderboardRedis.CalculateUserScore", zap.Error(err))
					nack = true
				}
				// done update score user :)

			default:
				nack = true
			}

			if nack {
				zap.L().Info(fmt.Sprintf("NAcking message from queue %s", "user-add-score"))

				_ = msg.Nack(false, nack)
			} else {
				zap.L().Info("Acking ")

				_ = msg.Ack(false)
			}

			zap.L().Info("No more messages to consume. Extiing.")

		}
	}()

	return nil

}

func decodeCorrectAnswer(b []byte) (domain.CorrectAnswer, error) {
	var res domain.CorrectAnswer
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		return domain.CorrectAnswer{}, domain.WrapErrorf(err, domain.ErrInternalServerError, "gob.Decode (decodeCorrectAnswer)")
	}
	return res, nil
}
