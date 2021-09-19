package main

import (
	"context"
	"github.com/Dementir/test/internal/api/routes"
	"github.com/Dementir/test/internal/store"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := initLog()
	defer logger.Sync() // flushes buffer, if any

	logger.Info("App start")

	db, err := sqlx.Open("pgx", "host=localhost port=5432 user=app password=pass database=job1 sslmode=disable")
	if err != nil {
		logger.Fatal(err)
	}

	pollRepo := store.New(db)

	dispatcher := routes.New(pollRepo, logger)
	httpServer := dispatcher.Init()

	go func() {
		logger.Info("starting HTTP server: listening on 10000")

		err := httpServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sgnl := make(chan os.Signal, 1)
	signal.Notify(sgnl,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	stop := <-sgnl

	if err = httpServer.Shutdown(context.Background()); err != nil {
		logger.Fatal("server shutdown with error: " + err.Error())
	}

	logger.Info("stopping", "signal", stop)
	logger.Info("waiting over jobs stopping")
}

func initLog() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := cfg.Build()

	return logger.Sugar()
}
