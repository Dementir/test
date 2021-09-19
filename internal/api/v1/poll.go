package v1

import (
	"encoding/json"
	"fmt"
	"github.com/Dementir/test/internal/customerror"
	"github.com/Dementir/test/internal/store"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
)

type getAnswerRequest struct {
	UserID int64 `json:"user,string"`
}

type getAnswerResponse struct {
	Question string   `json:"question"`
	Answers  []answer `json:"answers"`
}

type answer struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

type saveAnswerRequest struct {
	UserID int64  `json:"user,string"`
	Answer string `json:"answer"`
}

type saveAnswerResponse struct {
	Right bool `json:"right"`
}

type AnswerHandler struct {
	pollRepo store.PollRepository
	logger   *zap.SugaredLogger
}

func NewAnswerHandler(pollRepo store.PollRepository, logger *zap.SugaredLogger) *AnswerHandler {
	return &AnswerHandler{
		pollRepo: pollRepo,
		logger:   logger,
	}
}

func (ah *AnswerHandler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	req := new(getAnswerRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		errMsg := fmt.Sprintf("can not parse get answer request: %v", err)
		ah.logger.Error(errMsg)
		err = WithJSON(w, http.StatusBadRequest, &customerror.CustomError{
			ErrorCode:    300,
			ErrorMessage: errMsg,
		})

		return
	}

	question, err := ah.pollRepo.GetQuestion(r.Context(), req.UserID)
	if err != nil {
		errMsg := fmt.Sprintf("can not get question: %v", err)
		ah.logger.Error(errMsg)
		err = WithJSON(w, http.StatusBadRequest, &err)

		return
	}

	dbAnswers := make([]string, 0, len(question.WrongAnswers)+1)
	dbAnswers = append(dbAnswers, question.WrongAnswers...)
	dbAnswers = append(dbAnswers, question.RightAnswer)

	rand.Shuffle(len(dbAnswers), func(i, j int) {
		dbAnswers[i], dbAnswers[j] = dbAnswers[j], dbAnswers[i]
	})

	answers := make([]answer, 0, len(dbAnswers))
	for i, a := range dbAnswers {
		respAnswer := answer{
			Number: i + 1,
			Text:   a,
		}

		answers = append(answers, respAnswer)
	}

	err = WithJSON(w, http.StatusOK, &getAnswerResponse{
		Question: question.QuestionText,
		Answers:  answers,
	})
	if err != nil {
		ah.logger.Errorf("cannot get response: %v", err)

		return
	}
}

func (ah *AnswerHandler) IsAnswerRight(w http.ResponseWriter, r *http.Request) {
	req := new(saveAnswerRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		errMsg := fmt.Sprintf("can not parse get answer request: %v", err)
		ah.logger.Error(errMsg)
		err = WithJSON(w, http.StatusBadRequest, &customerror.CustomError{
			ErrorCode:    300,
			ErrorMessage: errMsg,
		})

		return
	}

	isRight, err := ah.pollRepo.IsAnswerRight(r.Context(), req.UserID, req.Answer)
	if err != nil {
		errMsg := fmt.Sprintf("can not get is answer: %v", err)
		ah.logger.Error(errMsg)
		err = WithJSON(w, http.StatusBadRequest, &err)

		return
	}

	err = WithJSON(w, http.StatusOK, &saveAnswerResponse{
		Right: isRight,
	})
	if err != nil {
		ah.logger.Errorf("cannot get response: %v", err)

		return
	}
}
