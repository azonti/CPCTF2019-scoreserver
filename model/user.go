package model

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
)

//User User
type User struct {
	ObjectID bson.ObjectId `bson:"_id"`
	Provider string        `bson:"provider"`
	ID       string        `bson:"id"`
	Token    string        `bson:"token"`
	UserInfo `bson:",inline"`
}

//UserInfo User Information
type UserInfo struct {
	Name              string `bson:"name"`
	IconURL           string `bson:"icon_url"`
	TwitterScreenName string `bson:"twitter_screen_name"`
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
			ObjectID: bson.NewObjectId(),
			Provider: provider,
			ID:       id,
		}
		if err := db.C("user").Insert(user); err != nil {
			return nil, fmt.Errorf("failed to insert a new user record: %v", err)
		}
	} else {
		user = &User{}
		if err := db.C("user").Find(bson.M{"provider": provider, "id": id}).One(user); err != nil {
			return nil, fmt.Errorf("failed to get the user record: %v", err)
		}
	}
	return user, nil
}

//SetToken Set a Token
func (user *User) SetToken() error {
	token := uuid.NewV4().String()
	if err := db.C("user").UpdateId(user.ObjectID, bson.M{"$set": bson.M{"token": token}}); err != nil {
		return fmt.Errorf("failed to update the user record: %v", err)
	}
	user.Token = token
	return nil
}
