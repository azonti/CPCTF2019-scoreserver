package model

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	webshell "git.trapti.tech/CPCTF2019/webshell/rpc"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/resty.v1"
)

//User an User Record
type User struct {
	ID                    string `gorm:"primary_key"`
	Token                 string
	TokenExpires          time.Time
	Name                  string
	IconURL               string
	TwitterScreenName     string
	IsAuthor              bool
	IsOnsite              bool
	Score                 int
	SolvedChallenges      []*Challenge `gorm:"many2many:user_solved_challenges;"`
	OpenedHints           []*Hint      `gorm:"many2many:user_opened_hints;"`
	FoundFlags            []*FoundFlag `gorm:"many2many:user_found_flags;"`
	WebShellPass          string
	LastSeenChallengeID   string
	LastSeenChallenge     *Challenge `gorm:"foreignkey:LastSeenChallengeID"`
	LastSolvedChallengeID string
	LastSolvedChallenge   *Challenge `gorm:"foreignkey:LastSolvedChallengeID"`
	LastSolvedTime        time.Time
	Votes                 []*Vote
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

//FoundFlag a FoundFlag Record
type FoundFlag struct {
	ID          int `gorm:"primary_key"`
	FlagID      string
	ChallengeID string
	Score       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

//Nobody a User Record which does Not Exist Actually
var Nobody = &User{
	ID: "nobody",
}
var appOnlyAuthConfig = map[string]*clientcredentials.Config{
	"twitter": {
		ClientID:     os.Getenv("TWITTER_CONSUMER_KEY"),
		ClientSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	},
}

//ErrUserNotFound an Error due to the User Not Found
var ErrUserNotFound = gorm.ErrRecordNotFound

//GetUsers Get All User Records
func GetUsers() ([]*User, error) {
	users := make([]*User, 0)
	if err := db.Preload("SolvedChallenges").Preload("OpenedHints").Preload("FoundFlags").Preload("Votes").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

//GetUserByID Get the User Record by their ID
func GetUserByID(id string, force bool) (*User, error) {
	user := &User{}
	err := db.Where(&User{ID: id}).Preload("SolvedChallenges").Preload("SolvedChallenges.Author").Preload("OpenedHints").Preload("FoundFlags").Preload("LastSeenChallenge").Preload("LastSolvedChallenge").Preload("Votes").First(user).Error
	if err == gorm.ErrRecordNotFound && force {
		name, iconURL, twitterScreenName, err := getUserInfo(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get the user's information: %v", err)
		}
		user = &User{
			ID:                id,
			Name:              name,
			IconURL:           iconURL,
			TwitterScreenName: twitterScreenName,
		}
		if err := db.Create(user).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

//GetUserByToken Get the User Record by their Token
func GetUserByToken(token string) (*User, error) {
	user := new(User)
	if err := db.Where(&User{Token: token}).Preload("SolvedChallenges").Preload("SolvedChallenges.Author").Preload("OpenedHints").Preload("FoundFlags").Preload("LastSeenChallenge").Preload("LastSolvedChallenge").Preload("Votes").First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//Delete Delete the User Record
func (user *User) Delete() error {
	return db.Delete(user).Error
}

//PutToken Put a Token
func (user *User) PutToken() error {
	token, tokenExpires := uuid.NewV4().String(), time.Now().Add(24*time.Hour)
	user.Token, user.TokenExpires = token, tokenExpires
	return db.Save(user).Error
}

//DeleteToken Delete the Token
func (user *User) DeleteToken() error {
	user.Token = ""
	return db.Save(user).Error
}

//MakeMeAuthor Make the User an Author
func (user *User) MakeMeAuthor() error {
	user.IsAuthor = true
	return db.Save(user).Error
}

//MakeMeOnsite Make the User Onsite
func (user *User) MakeMeOnsite() error {
	if err := user.RecreateWebShellContainer(); err != nil {
		return fmt.Errorf("failed to create the user's web shell container: %v", err)
	}
	user.IsOnsite = true
	return db.Save(user).Error
}

//OpenHint Open the Hint
func (user *User) OpenHint(id string) error {
	idSplit := strings.Split(id, ":")

	tx := db.Begin()
	hint := &Hint{}
	if err := tx.Where(&Hint{ID: id}).First(hint).Error; err != nil {
		tx.Rollback()
		return err
	}

	challenge := &Challenge{}
	if err := tx.Where(&Challenge{ID: idSplit[0]}).Preload("Flags").Find(&challenge).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(user).Association("OpenedHints").Append(hint).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.OpenedHints = append(user.OpenedHints, hint)

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

//RecreateWebShellContainer (Re)create the User's Web Shell Container
func (user *User) RecreateWebShellContainer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	webShellRes, err := webShellCli.New(ctx, &webshell.Request{
		Id:          user.ID,
		ScreenName:  map[bool]string{true: user.TwitterScreenName, false: user.ID}[user.TwitterScreenName != ""],
		DisplayName: user.Name,
	})
	if err != nil {
		return err
	}
	user.WebShellPass = webShellRes.GetPassword()
	return db.Save(user).Error
}

//SetLastSeenChallengeID Set the Challenge's ID which the User Saw Last
func (user *User) SetLastSeenChallengeID(challengeID string) error {
	user.LastSeenChallengeID = challengeID
	return db.Save(user).Error
}

func getUserInfo(id string) (string, string, string, error) {
	idSplit := strings.Split(id, "_")
	provider, rawID := idSplit[0], idSplit[1]
	httpClient := appOnlyAuthConfig[provider].Client(oauth2.NoContext)
	client := resty.New().SetTransport(httpClient.Transport)
	switch provider {
	case "twitter":
		data := &struct {
			Name            string `json:"name"`
			ScreenName      string `json:"screen_name"`
			ProfileImageURL string `json:"profile_image_url_https"`
		}{}
		if _, err := client.R().SetResult(data).Get("https://api.twitter.com/1.1/users/show.json?user_id=" + rawID); err != nil {
			return "", "", "", err
		}
		if data.Name == "" || data.ScreenName == "" || data.ProfileImageURL == "" {
			return "", "", "", fmt.Errorf("failed for unknown reason")
		}
		r := strings.NewReplacer("_normal", "")
		return data.Name, r.Replace(data.ProfileImageURL), data.ScreenName, nil
	}
	return "", "", "", ErrUnknownProvider
}
