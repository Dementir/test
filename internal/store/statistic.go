package store

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Statistic struct {
	Channel   string
	User      string
	Region    string
	ServiceID string
	Status    string
	Time      time.Time
}

type log struct {
	Stamp   time.Time `db:"stamp"`
	Channel string    `db:"channel"`
	Region  string    `db:"region"`
	Service string    `db:"service"`
	Status  string    `db:"status"`
	User    string    `db:"user"`
}

type Channel struct {
	ChannelID string `db:"channel_id"`
	Name      string `db:"name"`
}

type Region struct {
	RegionID string `db:"region_id"`
	Name     string `db:"name"`
}

type Service struct {
	ServiceID string `db:"service_id"`
	Name      string `db:"name"`
}

type Status struct {
	StatusID string `db:"status_id"`
	Name     string `db:"name"`
}

type statisticRepo struct {
	db *sqlx.DB
}

func NewStatisticRepository(db *sqlx.DB) StatisticRepository {
	return &statisticRepo{
		db: db,
	}
}

func (s *statisticRepo) Add(ctx context.Context, statistics []Statistic) error {

	channelQuery := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("channel").Columns("channel_id", "name")

	regionQuery := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("region").Columns("region_id", "name")

	serviceQuery := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("service").Columns("service_id", "name")

	statusQuery := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("status").Columns("status_id", "name")

	logQuery := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("log").Columns("stamp", "channel", "region", "service", "status", "\"user\"")

	i := 0
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, statistic := range statistics {
		if i >= 5000 {

			err = s.execQuery(ctx, tx, channelQuery)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			channelQuery = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("channel").Columns("channel_id", "name")

			err = s.execQuery(ctx, tx, regionQuery)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			regionQuery = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("region").Columns("region_id", "name")

			err = s.execQuery(ctx, tx, serviceQuery)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			serviceQuery = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("service").Columns("service_id", "name")

			err = s.execQuery(ctx, tx, statusQuery)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			statusQuery = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("status").Columns("status_id", "name")

			err = s.execQuery(ctx, tx, logQuery)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			logQuery = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert("log").Columns("stamp", "channel", "region", "service", "status", "\"user\"")

			i = 0
		}
		channelID := uuid.New().String()
		channel := Channel{
			ChannelID: channelID,
			Name:      statistic.Channel,
		}

		channelQuery = channelQuery.Values(channel.ChannelID, channel.Name)

		regionID := uuid.New().String()
		region := Region{
			RegionID: regionID,
			Name:     statistic.Channel,
		}

		regionQuery = regionQuery.Values(region.RegionID, region.Name)

		serviceID := uuid.New().String()
		service := Service{
			ServiceID: serviceID,
			Name:      statistic.Channel,
		}

		serviceQuery = serviceQuery.Values(service.ServiceID, service.Name)

		statusID := uuid.New().String()
		status := Status{
			StatusID: statusID,
			Name:     statistic.Channel,
		}

		statusQuery = statusQuery.Values(status.StatusID, status.Name)

		log := log{
			Stamp:   statistic.Time,
			Channel: channelID,
			Region:  regionID,
			Service: serviceID,
			Status:  statusID,
			User:    statistic.User,
		}

		logQuery = logQuery.Values(log.Stamp, log.Channel, log.Region, log.Service, log.Status, log.User)

		i++
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *statisticRepo) execQuery(ctx context.Context, tx *sqlx.Tx, query squirrel.InsertBuilder) error {
	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
