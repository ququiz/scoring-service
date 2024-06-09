package mongodb

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type QuizRepository struct {
	db *mongo.Database
}

func NewQuizRepository(db *mongo.Database) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) UpdateParticipantScore(ctx context.Context, quizID, participantUserID string, finalScore uint64) error {
	coll := r.db.Collection("base_quiz")

	filterQuiz := bson.D{{"_id", quizID}, {"$participants.user_id", participantUserID}}
	updateParticipantScore := bson.D{{"$set", bson.D{{"$participants.final_score", finalScore}}}}
	_, err := coll.UpdateOne(context.Background(), filterQuiz, updateParticipantScore)
	if err != nil {
		zap.L().Error("coll.UpdateOne (UpdateParticipantScore) (UpdateParticipantScore) ")
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	return nil
}
