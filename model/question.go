package model

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

//Question a Question Record
type Question struct {
	ID           string `gorm:"primary_key"`
	QuestionerID string
	Questioner   *User `gorm:"foreignkey:QuestionerID"`
	AnswererID   string
	Answerer     *User `gorm:"foreignkey:AnswererID"`
	Query        string
	Answer       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

//ErrQuestionNotFound an Error due to the Question Not Found
var ErrQuestionNotFound = gorm.ErrRecordNotFound

//GetQuestions Get All Question Records
func GetQuestions() ([]*Question, error) {
	questions := make([]*Question, 0)
	if err := db.Preload("Questioner").Preload("Answerer").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

//GetQuestionByID Get the Question Record by its ID
func GetQuestionByID(id string) (*Question, error) {
	question := new(Question)
	if err := db.Where(&Question{ID: id}).Preload("Questioner").Preload("Answerer").First(question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

//NewQuestion Make a New Question Record
func NewQuestion(questionerID string, query string) (*Question, error) {
	id := uuid.NewV4().String()
	question := &Question{
		ID:           id,
		QuestionerID: questionerID,
		Query:        query,
	}
	if err := db.Create(question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

//Update Update the Question Record
func (question *Question) Update(questionerID string, answererID string, query string, answer string) error {
	question.QuestionerID, question.AnswererID, question.Query, question.Answer = questionerID, answererID, query, answer
	return db.Save(question).Error
}
