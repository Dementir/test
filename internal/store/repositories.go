package store

import "context"

type StatisticRepository interface {
	Add(ctx context.Context, statistic []Statistic) error
}

type PollRepository interface {
	GetQuestion(ctx context.Context, userID int64) (*Question, error)
	IsAnswerRight(ctx context.Context, userID int64, answer string) (bool, error)
}
