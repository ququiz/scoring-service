package mongodb

import (
	"context"
	"fmt"
	"ququiz/lintang/scoring-service/biz/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	quizIDObjectID, err := primitive.ObjectIDFromHex(quizID)
	filterQuiz := bson.D{{"_id", quizIDObjectID}, {"participants.user_id", participantUserID}}

	zap.L().Debug(fmt.Sprintf(`userID: %s , final_score: %s`, participantUserID, finalScore))
	updateParticipantScore := bson.D{{"$set", bson.D{{"participants.$.final_score", finalScore}}}}

	_, err = coll.UpdateOne(context.Background(), filterQuiz, updateParticipantScore)
	if err != nil {
		zap.L().Error("coll.FindOneAndUpdate (UpdateParticipantScore) (UpdateParticipantScore) ", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	return nil
}
