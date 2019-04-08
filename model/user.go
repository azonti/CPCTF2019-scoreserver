package model

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/resty.v1"
)

//User an User Record
type User struct {
	ObjectID              bson.ObjectId `bson:"_id"`
	ID                    string        `bson:"id"`
	Token                 string        `bson:"token"`
	TokenExpires          time.Time     `bson:"token_expires"`
	Name                  string        `bson:"name"`
	IconURL               string        `bson:"icon_url"`
	TwitterScreenName     string        `bson:"twitter_screen_name"`
	IsAuthor              bool          `bson:"is_author"`
	IsOnsite              bool          `bson:"is_onsite"`
	OpenedHintIDs         []string      `bson:"opened_hint_ids"`
	WebShellPass          string        `bson:"web_shell_pass"`
	LastSeenChallengeID   string        `bson:"last_seen_challenge_id"`
	LastSolvedChallengeID string        `bson:"last_solved_challenge_id"`
	LastSolvedTime        time.Time     `bson:"last_solved_time"`
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
var ErrUserNotFound = fmt.Errorf("the user not found")

//GetUsers Get All User Records
func GetUsers() ([]*User, error) {
	var users []*User
	if err := db.C("user").Find(nil).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

//GetUserByID Get the User Record by their ID
func GetUserByID(id string, force bool) (*User, error) {
	if force {
		n, err := db.C("user").Find(bson.M{"id": id}).Count()
		if err != nil {
			return nil, fmt.Errorf("failed to check the user record existence: %v", err)
		}
		if n == 0 {
			name, iconURL, twitterScreenName, err := getUserInfo(id)
			if err != nil {
				return nil, fmt.Errorf("failed to get the user's information: %v", err)
			}
			user := &User{
				ObjectID:          bson.NewObjectId(),
				ID:                id,
				Name:              name,
				IconURL:           iconURL,
				TwitterScreenName: twitterScreenName,
			}
			if err := db.C("user").Insert(user); err != nil {
				return nil, fmt.Errorf("failed to insert a new user record: %v", err)
			}
			return user, nil
		}
	}
	user := &User{}
	if err := db.C("user").Find(bson.M{"id": id}).One(user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

//GetUserByToken Get the User Record by their Token
func GetUserByToken(token string) (*User, error) {
	user := &User{}
	if err := db.C("user").Find(bson.M{"token": token, "token_expires": bson.M{"$gte": time.Now()}}).One(user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

//Delete Delete the User Record
func (user *User) Delete() error {
	return db.C("user").RemoveId(user.ObjectID)
}

//SetToken Set a Token
func (user *User) SetToken() error {
	token, tokenExpires := uuid.NewV4().String(), time.Now().Add(24*time.Hour)
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"token": token, "token_expires": tokenExpires}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.Token, user.TokenExpires = token, tokenExpires
	return nil
}

//RemoveToken Remove the Token
func (user *User) RemoveToken() error {
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"token": ""}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.Token = ""
	return nil
}

//MakeMeAuthor Make the User a Author
func (user *User) MakeMeAuthor() error {
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"is_author": true}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.IsAuthor = true
	return nil
}

//MakeMeOnsite Make the User Onsite
func (user *User) MakeMeOnsite() error {
	if err := user.RecreateWebShellContainer(); err != nil {
		return fmt.Errorf("failed to create the user's web shell container: %v", err)
	}
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"is_onsite": true}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.IsOnsite = true
	return nil
}

//OpenHint Open the Hint
func (user *User) OpenHints(id []string) error {
	newOpenedHintIDs := append(user.OpenedHintIDs, id...)
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"opened_hint_ids": newOpenedHintIDs}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.OpenedHintIDs = newOpenedHintIDs
	return nil
}

//RecreateWebShellContainer (Re)create the User's Web Shell Container
func (user *User) RecreateWebShellContainer() error {
	return fmt.Errorf("not implementated")
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// webShellRes, err := webShellCli.Create(ctx, &webshell.Request{
	// 	Id:          user.ID,
	// 	ScreenName:  map[bool]string{true: user.TwitterScreenName, false: user.ID}[user.TwitterScreenName != ""],
	// 	DisplayName: user.Name,
	// })
	// if err != nil {
	// 	return err
	// }
	// if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"web_shell_pass": webShellRes.GetPassword()}}); err != nil {
	// 	return fmt.Errorf("failed to update the user record: %v", err)
	// }
	// user.WebShellPass = webShellRes.GetPassword()
	// return nil
}

//GetSolvedChallenges Get the Challenges which the User solved
func (user *User) GetSolvedChallenges() ([]*Challenge, error) {
	pipe := db.C("challenge").Pipe([]bson.M{
		{"$match": bson.M{"who_solved_ids": user.ID}},
	})
	var challenges []*Challenge
	if err := pipe.All(&challenges); err != nil {
		return nil, err
	}
	return challenges, nil
}

var scoreCache = map[string]int{}

//GetScore Get the User's Score
func (user *User) GetScore() (int, error) {
	if v, ok := scoreCache[user.ID]; ok {
		return v, nil
	}
	rawScorePipe, penaltyPipe := db.C("challenge").Pipe([]bson.M{
		{"$project": bson.M{"score": 1, "who_pointed_ids": 1}},
		{"$match": bson.M{"who_pointed_ids": user.ID}},
		{"$group": bson.M{"_id": "score", "score": bson.M{"$sum": "$score"}}},
	}), db.C("challenge").Pipe([]bson.M{
		{"$project": bson.M{"hints": 1, "who_pointed_ids": 1}},
		{"$match": bson.M{"who_pointed_ids": user.ID}},
		{"$unwind": "$hints"},
		{"$match": bson.M{"hints.id": bson.M{"$in": user.OpenedHintIDs}}},
		{"$group": bson.M{"_id": "penalty", "penalty": bson.M{"$sum": "$hints.penalty"}}},
	})
	rawScore, penalty := &struct {
		ObjectID string `bson:"_id"`
		Score    int    `bson:"score"`
	}{}, &struct {
		ObjectID string `bson:"_id"`
		Penalty  int    `bson:"penalty"`
	}{}
	if err := rawScorePipe.One(rawScore); err != nil {
		if err != mgo.ErrNotFound {
			return 0, fmt.Errorf("failed to get the user's raw score: %v", err)
		}
		rawScore.Score = 0
	}
	if err := penaltyPipe.One(penalty); err != nil {
		if err != mgo.ErrNotFound {
			return 0, fmt.Errorf("failed to get the user's penalty: %v", err)
		}
		penalty.Penalty = 0
	}
	scoreCache[user.ID] = rawScore.Score - penalty.Penalty
	return scoreCache[user.ID], nil
}

func (user *User) setLastSolvedChallengeID(challengeID string) error {
	now := time.Now()
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"last_solved_challenge_id": challengeID, "last_solved_time": now}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.LastSolvedChallengeID, user.LastSolvedTime = challengeID, now
	return nil
}

//SetLastSeenChallengeID Set the Challenge's ID which the User Saw Last
func (user *User) SetLastSeenChallengeID(challengeID string) error {
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"last_seen_challenge_id": challengeID}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.LastSeenChallengeID = challengeID
	return nil
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
