package model

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"strconv"
)

//Challenge a Challenge Record
type Challenge struct {
	ObjectID     bson.ObjectId `bson:"_id"`
	ID           string        `bson:"id"`
	AuthorID     string        `bson:"author_id"`
	Score        int           `bson:"score"`
	Caption      string        `bson:"caption"`
	Hints        []*Hint       `bson:"hints"`
	Flag         string        `bson:"flag"`
	Answer       string        `bson:"answer"`
	WhoSolvedIDs []string      `bson:"who_solved_ids"`
}

//Hint a Hint Record
type Hint struct {
	ID      string `bson:"id"`
	Caption string `bson:"caption"`
	Penalty int    `bson:"penalty"`
}

//ErrChallengeNotFound an Error due to the Challenge Not Found
var ErrChallengeNotFound = fmt.Errorf("the challenge not found")

//GetChallenges Get All Challenge Records
func GetChallenges() ([]*Challenge, error) {
	var challenges []*Challenge
	if err := db.C("challenge").Find(nil).All(&challenges); err != nil {
		return nil, err
	}
	return challenges, nil
}

//GetChallengeByID Get the Challenge Record by its ID
func GetChallengeByID(id string) (*Challenge, error) {
	challenge := &Challenge{}
	if err := db.C("challenge").Find(bson.M{"id": id}).One(challenge); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}
	return challenge, nil
}

//NewChallenge Make a New Challenge Record
func NewChallenge(authorID string, score int, caption string, captions []string, penalties []int, flag string, answer string) (*Challenge, error) {
	challenge := &Challenge{
		ObjectID: bson.NewObjectId(),
		ID:       uuid.NewV4().String(),
		AuthorID: authorID,
		Score:    score,
		Caption:  caption,
		Hints:    make([]*Hint, len(captions)),
		Flag:     flag,
		Answer:   answer,
	}
	for i := 0; i < len(captions); i++ {
		challenge.Hints[i] = &Hint{
			ID:      challenge.ID + ":" + strconv.Itoa(i),
			Caption: captions[i],
			Penalty: penalties[i],
		}
	}
	if err := db.C("challenge").Insert(challenge); err != nil {
		return nil, fmt.Errorf("failed to insert a new challenge record: %v", err)
	}
	return challenge, nil
}

//Delete Delete the Challenge Record
func (challenge *Challenge) Delete() error {
	return db.C("challenge").RemoveId(challenge.ObjectID)
}

//Update Update the Challenge Record
func (challenge *Challenge) Update(authorID string, score int, caption string, captions []string, penalties []int, flag string, answer string) error {
	hints := make([]bson.M, len(captions))
	for i := 0; i < len(captions); i++ {
		hints[i] = bson.M{"id": challenge.ID + ":" + strconv.Itoa(i), "caption": captions[i], "penalty": penalties[i]}
	}
	if err := db.C("challenge").UpdateId(challenge.ObjectID, bson.M{"$set": bson.M{"author_id": authorID, "score": score, "caption": caption, "hints": hints, "flag": flag, "answer": answer}}); err != nil {
		return err
	}
	challenge.AuthorID, challenge.Score, challenge.Caption, challenge.Hints, challenge.Flag, challenge.Answer = authorID, score, caption, make([]*Hint, len(captions)), flag, answer
	for i := 0; i < len(captions); i++ {
		challenge.Hints[i] = &Hint{
			ID:      challenge.ID + ":" + strconv.Itoa(i),
			Caption: captions[i],
			Penalty: penalties[i],
		}
	}
	return nil
}

//AddWhoSolved Add the User to the List of Who Solved
func (challenge *Challenge) AddWhoSolved(user *User) error {
	newWhoSolvedIDs := append(challenge.WhoSolvedIDs, user.ID)
	if err := db.C("challenge").UpdateId(challenge.ObjectID, bson.M{"$set": bson.M{"who_solved_ids": newWhoSolvedIDs}}); err != nil {
		return fmt.Errorf("failed to update the challenge record: %v", err)
	}
	challenge.WhoSolvedIDs = newWhoSolvedIDs
	return nil
}
