package routes

import (
	v1 "github.com/Dementir/test/internal/api/v1"
	"github.com/Dementir/test/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Dispatcher struct {
	pollRepo store.PollRepository
	logger   *zap.SugaredLogger
}

func New(pollRepo store.PollRepository, logger *zap.SugaredLogger) *Dispatcher {
	return &Dispatcher{
		pollRepo: pollRepo,
		logger:   logger,
	}
}

func (d *Dispatcher) Init() *http.Server {
	r := chi.NewRouter()

	answer := v1.NewAnswerHandler(d.pollRepo, d.logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/get_question", answer.GetAnswer)
		r.Post("/answer", answer.IsAnswerRight)
	})

	return &http.Server{
		Addr:         ":10000",
		Handler:      r,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
}
