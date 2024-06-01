package rediscache

import (
	"context"
	"fmt"
	"ququiz/lintang/scoring-service/biz/dal/domain"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LeaderboardRedis struct {
	cli *redis.Client
}

func NewLeaderboardRedis(cli *redis.Client) *LeaderboardRedis {
	return &LeaderboardRedis{cli}
}

func (c *LeaderboardRedis) CalculateUserScore(ctx context.Context, weight uint64, userID string, quizID string) error {
	key := "leaderboard" + quizID
	_, err := c.cli.ZRevRank(ctx, key, userID).Result()
	if err == redis.Nil {
		// user belum ada di leaderboard, insert ke leadeer board dg score 0 + weight
		_, err := c.cli.ZAdd(ctx, key, redis.Z{Score: float64(weight), Member: userID}).Result()
		if err != nil {
			zap.L().Error(fmt.Sprintf("error (c.cli.ZAdd) (CalculateUserScore) (LeaderboardRedis)", zap.Error(err)))
			return err
		}

		return nil
	}
	if err != nil {
		zap.L().Error(fmt.Sprintf("error (c.cli.ZRevRank) (CalculateUserScore) (LeaderboardRedis)", zap.Error(err)))
		return err
	}

	// kalau sebelumnya userIDnya ada di redis , increment by weight (lastscore + weight)
	_, err = c.cli.ZIncrBy(ctx, key, float64(weight), userID).Result()
	if err != nil {
		zap.L().Error("c.cli.ZIncrBy (CalculateUserScore) (LeaderboardRedis)", zap.Error(err))
		return err
	}

	return nil
}

func (c *LeaderboardRedis) GetTopLeaderBoard(ctx context.Context, quizID string) ([]domain.RedisLeaderBoard, error) {
	key := "leaderboard" + quizID
	l, err := c.cli.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err == redis.Nil {
		zap.L().Debug(fmt.Sprintf("leaderboard %s belum ada di sorted sets", key))
		return []domain.RedisLeaderBoard{}, domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf(fmt.Sprintf("leaderboard %s belum ada di sorted sets", key)))
	}
	if err != nil {
		zap.L().Error("c.cli.ZRevRange (GetTopLeaderBoard) (LeaderBoardRedis)", zap.Error(err))
		return []domain.RedisLeaderBoard{}, err
	}

	var leaderboards []domain.RedisLeaderBoard
	for i := 0; i < len(l); i++ {
		leaderboards = append(leaderboards, domain.RedisLeaderBoard{
			UserID:   l[i].Member.(string),
			Score:    uint64(l[i].Score),
			Position: uint64(i + 1),
		})
	}

	return leaderboards, nil
}
