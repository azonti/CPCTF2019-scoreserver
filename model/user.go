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

//User User
type User struct {
	ObjectID          bson.ObjectId `bson:"_id,omitempty"`
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

//GetUserByID Get the User by the Identity Provider's Name and the User ID
func GetUserByID(provider string, id string, force bool) (*User, error) {
	n, err := db.C("user").Find(bson.M{"provider": provider, "id": id}).Count()
	if err != nil {
		return nil, fmt.Errorf("failed to check the user record existence: %v", err)
	}
	var user *User
	if n == 0 && force {
		user = &User{
			Provider: provider,
			ID:       id,
		}
		if err := db.C("user").Insert(user); err != nil {
			return nil, fmt.Errorf("failed to insert the user record: %v", err)
		}
		go func() {
			if err := user.initUserInfo(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to init the user (%s:%s) info: %v", user.Provider, user.ID, err)
			}
			return
		}()
	} else {
		user = &User{}
		if err := db.C("user").Find(bson.M{"provider": provider, "id": id}).One(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

//GetUserByToken Get the User by the Token
func GetUserByToken(token string) (*User, error) {
	user := &User{}
	if err := db.C("user").Find(bson.M{"token": token}).One(user); err != nil {
		return nil, fmt.Errorf("failed to get the user record: %v", err)
	}
	return user, nil
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

func (user *User) initUserInfo() error {
	httpClient := appOnlyAuthConfig[user.Provider].Client(oauth2.NoContext)
	client := resty.New().SetTransport(httpClient.Transport)
	switch user.Provider {
	case "twitter":
		data := &struct {
			Name            string `json:"name"`
			ScreenName      string `json:"screen_name"`
			ProfileImageURL string `json:"profile_image_url"`
		}{}
		if _, err := client.R().SetResult(data).Get("https://api.twitter.com/1.1/users/show.json?user_id=" + user.ID); err != nil {
			return fmt.Errorf("failed to get the user info: %v", err)
		}
		if data.Name == "" || data.ScreenName == "" || data.ProfileImageURL == "" {
			return fmt.Errorf("failed for unknown reason")
		}
		if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"name": data.Name, "icon_url": data.ProfileImageURL, "twitter_screen_name": data.ScreenName}}); err != nil {
			return fmt.Errorf("failed to update the user record: %v", err)
		}
		user.Name, user.IconURL, user.TwitterScreenName = data.Name, data.ProfileImageURL, data.ScreenName
	}
	return nil
}
