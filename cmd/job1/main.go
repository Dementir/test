package main

import (
	"context"
	"flag"
	"github.com/Dementir/test/internal/logparser"
	"github.com/Dementir/test/internal/store"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func main() {
	logPath := flag.String("log", "logfile.log", "set log path")
	flag.Parse()

	db, err := sqlx.Open("pgx", "host=localhost port=5432 user=app password=pass database=job1 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	statRepo := store.NewStatisticRepository(db)

	stats, err := logparser.LogParse(*logPath)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*20)
	defer cancel()

	err = statRepo.Add(ctx, stats)
	if err != nil {
		log.Fatalln(err)
	}
}
