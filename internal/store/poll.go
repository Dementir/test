package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/Dementir/test/internal/customerror"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type Question struct {
	QuestionID   int64          `db:"question_id"`
	QuestionText string         `db:"question_text"`
	RightAnswer  string         `db:"right_answer"`
	WrongAnswers pq.StringArray `db:"wrong_answer"`
}

type User struct {
	UserID int64  `db:"user_id"`
	Name   string `db:"name"`
}

type Game struct {
	UserID      int64  `db:"user"`
	QuestionID  int64  `db:"question"`
	RightAnswer string `db:"right_answer"`
	Answered    bool   `db:"answered"`
}

type Rating struct {
	UserID      int64     `db:"user"`
	RightAnswer int64     `db:"right_answers"`
	AnswerTime  time.Time `db:"answer_time"`
}

type pollRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) PollRepository {
	return &pollRepo{
		db: db,
	}
}

func (p *pollRepo) GetQuestion(ctx context.Context, userID int64) (*Question, error) {
	question, err := p.getQuestion(ctx, userID)
	if err != nil {
		return nil, &customerror.CustomError{
			ErrorCode:    100,
			ErrorMessage: fmt.Sprintf("get question error: %v", err),
		}
	}

	err = p.saveGameState(ctx, userID, *question, false)
	if err != nil {
		return nil, &customerror.CustomError{
			ErrorCode:    101,
			ErrorMessage: fmt.Sprintf("save game state error: %v", err),
		}
	}

	return question, nil
}

func (p *pollRepo) getQuestion(ctx context.Context, userID int64) (*Question, error) {
	const query = `SELECT question_id
	, question_text
	, right_answer
	, wrong_answer
FROM question ORDER BY random() LIMIT 1;`

	question := new(Question)
	err := p.db.GetContext(ctx, question, query)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (p *pollRepo) saveGameState(ctx context.Context, userID int64, question Question, isAnswered bool) error {
	const query = `insert into game(
	"user"
	, question
	, right_answer
	, answered
) values (
	:user
	, :question
	, :right_answer
	, :answered);`

	game := &Game{
		UserID:      userID,
		QuestionID:  question.QuestionID,
		RightAnswer: question.RightAnswer,
		Answered:    isAnswered,
	}

	_, err := p.db.NamedExecContext(ctx, query, game)
	if err != nil {
		return err
	}

	return nil
}

func (p *pollRepo) isAnswered(ctx context.Context, userID int64) (bool, error) {
	const query = `select answered from game where "user" = $1;`

	row := p.db.QueryRowContext(ctx, query, userID)
	err := row.Err()
	if err != nil {
		return false, fmt.Errorf("get is answered error: %v", err)
	}

	var isAnswered bool
	err = row.Scan(&isAnswered)
	if err != nil {
		return false, fmt.Errorf("scan is answered error: %v", err)
	}

	return isAnswered, nil
}

func (p *pollRepo) IsAnswerRight(ctx context.Context, userID int64, answer string) (bool, error) {
	isAnswered, err := p.isAnswered(ctx, userID)
	if err != nil {
		return false, err
	}

	if isAnswered {
		return false, &customerror.CustomError{
			ErrorCode:    102,
			ErrorMessage: fmt.Sprintf("get question error: %v", errors.New("user has already answered")),
		}
	}

	question, err := p.getQuestion(ctx, userID)
	if err != nil {
		return false, &customerror.CustomError{
			ErrorCode:    100,
			ErrorMessage: fmt.Sprintf("get question error: %v", err),
		}
	}

	err = p.saveGameState(ctx, userID, *question, true)
	if err != nil {
		return false, &customerror.CustomError{
			ErrorCode:    101,
			ErrorMessage: fmt.Sprintf("save game state error: %v", err),
		}
	}

	return question.RightAnswer == answer, nil
}
