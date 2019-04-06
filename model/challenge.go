package model

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)

//Challenge a Challenge Record
type Challenge struct {
	ID            string `gorm:"primary_key"`
	Genre         string
	Name          string
	AuthorID      string
	Author        *User `gorm:"foreignkey:AuthorID"`
	Score         int
	Caption       string
	Hints         []*Hint
	Flag          string
	Answer        string
	WhoSolved     []*User `gorm:"many2many:user_solved_challenges;"`
	WhoChallenged []*User `gorm:"many2many:user_challenged_challenges;"`
	Votes         []*Vote
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

//Hint a Hint Record
type Hint struct {
	ID          string `gorm:"primary_key"`
	ChallengeID string
	Caption     string
	Penalty     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

//Vote a Vote Record
type Vote struct {
	gorm.Model
	ChallengeID string
	UserID      string
	Vote        string
}

//ErrChallengeNotFound an Error due to the Challenge Not Found
var ErrChallengeNotFound = gorm.ErrRecordNotFound

//GetChallenges Get All Challenge Records
func GetChallenges() ([]*Challenge, error) {
	challenges := make([]*Challenge, 0)
	if err := db.Preload("Author").Preload("Hints").Preload("WhoSolved").Preload("WhoChallenged").Preload("Votes").Find(&challenges).Error; err != nil {
		return nil, err
	}
	return challenges, nil
}

//GetChallengeByID Get the Challenge Record by its ID
func GetChallengeByID(id string) (*Challenge, error) {
	challenge := &Challenge{}
	if err := db.Where(&Challenge{ID: id}).Preload("Author").Preload("Hints").Preload("WhoSolved").Preload("WhoChallenged").Preload("Votes").First(challenge).Error; err != nil {
		return nil, err
	}
	return challenge, nil
}

//NewChallenge Make a New Challenge Record
func NewChallenge(genre string, name string, authorID string, score int, caption string, captions []string, penalties []int, flag string, answer string) (*Challenge, error) {
	id := uuid.NewV4().String()
	author, err := GetUserByID(authorID, false)
	if err != nil {
		return nil, err
	}
	hints := make([]*Hint, len(captions))
	for i := 0; i < len(captions); i++ {
		hints[i] = &Hint{
			ID:      id + ":" + strconv.Itoa(i),
			Caption: captions[i],
			Penalty: penalties[i],
		}
	}
	challenge := &Challenge{
		ID:      id,
		Genre:   genre,
		Name:    name,
		Author:  author,
		Score:   score,
		Caption: caption,
		Hints:   hints,
		Flag:    flag,
		Answer:  answer,
	}
	if err := db.Set("gorm:save_associations", true).Create(challenge).Error; err != nil {
		return nil, err
	}
	return challenge, nil
}

//Delete Delete the Challenge Record
func (challenge *Challenge) Delete() error {
	return db.Delete(challenge).Error
}

//Update Update the Challenge Record
func (challenge *Challenge) Update(genre string, name string, authorID string, score int, caption string, captions []string, penalties []int, flag string, answer string) error {
	author, err := GetUserByID(authorID, false)
	if err != nil {
		return err
	}
	hints := make([]*Hint, len(captions))
	for i := 0; i < len(captions); i++ {
		hints[i] = &Hint{
			ID:      challenge.ID + ":" + strconv.Itoa(i),
			Caption: captions[i],
			Penalty: penalties[i],
		}
	}
	challenge.Genre, challenge.Name, challenge.Author, challenge.Score, challenge.Caption, challenge.Hints, challenge.Flag, challenge.Answer = genre, name, author, score, caption, hints, flag, answer
	return db.Set("gorm:save_associations", true).Save(challenge).Error
}

//AddWhoSolved Add the User to the List of Who Solved
func (challenge *Challenge) AddWhoSolved(user *User) error {
	hints := make([]*Hint, 0)
	if err := db.Where(&Hint{ChallengeID: challenge.ID}).Model(user).Association("OpenedHints").Find(&hints).Error; err != nil {
		return err
	}
	scoreDelta := challenge.Score
	for _, hint := range hints {
		scoreDelta -= hint.Penalty
	}
	user.Score += scoreDelta
	user.LastSolvedChallenge = challenge
	user.LastSolvedTime = time.Now()
	tx := db.Begin()
	if err := tx.Model(challenge).Association("WhoSolved").Append(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//AddWhoChallenged Add the User to the List of Who Challenged
func (challenge *Challenge) AddWhoChallenged(user *User) error {
	return db.Model(challenge).Association("WhoChallenged").Append(user).Error
}

//GetVote Get the User's Vote for the Challenge
func (challenge *Challenge) GetVote(userID string) (string, error) {
	vote := &Vote{}
	if err := db.Where(&Vote{ChallengeID: challenge.ID, UserID: userID}).First(&vote).Error; err != nil {
		return "", err
	}
	return vote.Vote, nil
}

//PutVote Put the User's Vote for the Challenge
func (challenge *Challenge) PutVote(userID string, vote string) error {
	_vote := &Vote{
		ChallengeID: challenge.ID,
		UserID:      userID,
		Vote:        vote,
	}
	return db.Create(_vote).Error
}
