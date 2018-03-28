package model

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/resty.v1"
	"os"
	"time"
)

//User an User Record
type User struct {
	ObjectID          bson.ObjectId `bson:"_id"`
	Provider          string        `bson:"provider"`
	ID                string        `bson:"id"`
	Token             string        `bson:"token"`
	TokenExpires      time.Time     `bson:"token_expires"`
	Name              string        `bson:"name"`
	IconURL           string        `bson:"icon_url"`
	TwitterScreenName string        `bson:"twitter_screen_name"`
}

var appOnlyAuthConfig = map[string]*clientcredentials.Config{
	"twitter": {
		ClientID:     os.Getenv("TWITTER_CONSUMER_KEY"),
		ClientSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	},
}

//GetUserByID Get the User Record by their ID
func GetUserByID(provider string, id string, force bool) (*User, error) {
	if force {
		n, err := db.C("user").Find(bson.M{"provider": provider, "id": id}).Count()
		if err != nil {
			return nil, fmt.Errorf("failed to check the user record existence: %v", err)
		}
		if n == 0 {
			name, iconURL, twitterScreenName, err := getUserInfo(provider, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get the user's information: %v", err)
			}
			user := &User{
				ObjectID:          bson.NewObjectId(),
				Provider:          provider,
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
	if err := db.C("user").Find(bson.M{"provider": provider, "id": id}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

//GetUserByToken Get the User Record by their Token
func GetUserByToken(token string) (*User, error) {
	user := &User{}
	if err := db.C("user").Find(bson.M{"token": token}).One(user); err != nil {
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

func getUserInfo(provider string, id string) (string, string, string, error) {
	httpClient := appOnlyAuthConfig[provider].Client(oauth2.NoContext)
	client := resty.New().SetTransport(httpClient.Transport)
	switch provider {
	case "twitter":
		data := &struct {
			Name            string `json:"name"`
			ScreenName      string `json:"screen_name"`
			ProfileImageURL string `json:"profile_image_url"`
		}{}
		if _, err := client.R().SetResult(data).Get("https://api.twitter.com/1.1/users/show.json?user_id=" + id); err != nil {
			return "", "", "", err
		}
		if data.Name == "" || data.ScreenName == "" || data.ProfileImageURL == "" {
			return "", "", "", fmt.Errorf("failed for unknown reason")
		}
		return data.Name, data.ProfileImageURL, data.ScreenName, nil
	}
	return "", "", "", fmt.Errorf("an unknown provider")
}
