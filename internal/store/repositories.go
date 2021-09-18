package store

import "context"

type StatisticRepository interface {
	Add(ctx context.Context, statistic []Statistic) error
}
