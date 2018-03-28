package model

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/resty.v1"
	"os"
	"strings"
	"time"
)

//User an User Record
type User struct {
	ObjectID          bson.ObjectId `bson:"_id"`
	ID                string        `bson:"id"`
	Token             string        `bson:"token"`
	TokenExpires      time.Time     `bson:"token_expires"`
	Name              string        `bson:"name"`
	IconURL           string        `bson:"icon_url"`
	TwitterScreenName string        `bson:"twitter_screen_name"`
	IsAuthor          bool          `bson:"is_author"`
	OpenedHintIDs     []string      `bson:"opened_hint_ids"`
}

//Nobody a User Record which does Not Exist Actually
var Nobody = &User{
	ID:      "nobody:0",
	Name:    "Nobody",
	IconURL: os.Getenv("NOBODY_ICON_URL"),
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

func getUserInfo(id string) (string, string, string, error) {
	idSplit := strings.Split(id, ":")
	provider, rawID := idSplit[0], idSplit[1]
	httpClient := appOnlyAuthConfig[provider].Client(oauth2.NoContext)
	client := resty.New().SetTransport(httpClient.Transport)
	switch provider {
	case "twitter":
		data := &struct {
			Name            string `json:"name"`
			ScreenName      string `json:"screen_name"`
			ProfileImageURL string `json:"profile_image_url"`
		}{}
		if _, err := client.R().SetResult(data).Get("https://api.twitter.com/1.1/users/show.json?user_id=" + rawID); err != nil {
			return "", "", "", err
		}
		if data.Name == "" || data.ScreenName == "" || data.ProfileImageURL == "" {
			return "", "", "", fmt.Errorf("failed for unknown reason")
		}
		return data.Name, data.ProfileImageURL, data.ScreenName, nil
	}
	return "", "", "", ErrUnknownProvider
}
