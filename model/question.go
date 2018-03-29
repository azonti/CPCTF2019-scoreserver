package model

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
)

//Question a Question Record
type Question struct {
	ObjectID     bson.ObjectId `bson:"_id"`
	ID           string        `bson:"id"`
	QuestionerID string        `bson:"questioner_id"`
	AnswererID   string        `bson:"answerer_id"`
	Query        string        `bson:"query"`
	Answer       string        `bson:"answer"`
}

//ErrQuestionNotFound an Error due to the Question Not Found
var ErrQuestionNotFound = fmt.Errorf("the question not found")

//GetQuestions Get All Question Records
func GetQuestions() ([]*Question, error) {
	var questions []*Question
	if err := db.C("question").Find(nil).All(&questions); err != nil {
		return nil, err
	}
	return questions, nil
}

//GetQuestionByID Get the Question Record by its ID
func GetQuestionByID(id string) (*Question, error) {
	question := &Question{}
	if err := db.C("question").Find(bson.M{"id": id}).One(question); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	return question, nil
}

//NewQuestion Make a New Question Record
func NewQuestion(questionerID string, query string) (*Question, error) {
	question := &Question{
		ObjectID:     bson.NewObjectId(),
		ID:           uuid.NewV4().String(),
		QuestionerID: questionerID,
		Query:        query,
	}
	if err := db.C("question").Insert(question); err != nil {
		return nil, fmt.Errorf("failed to insert a new question record: %v", err)
	}
	return question, nil
}

//Update Update the Question Record
func (question *Question) Update(questionerID string, answererID string, query string, answer string) error {
	return db.C("question").UpdateId(question.ObjectID, bson.M{"$set": bson.M{"questioner_id": questionerID, "answerer_id": answererID, "query": query, "answer": answer}})
}
