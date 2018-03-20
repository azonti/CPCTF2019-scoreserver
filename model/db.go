package model

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
)

var session *mgo.Session

//DB Database
var DB *mgo.Database

//InitDB Initialize Database
func InitDB() error {
	var err error
	session, err = mgo.Dial(os.Getenv("MONGODB_URL"))
	if err != nil {
		return fmt.Errorf("failed to establish DB session")
	}
	DB = session.DB("")
	return nil
}

//TermDB Terminate Database
func TermDB() {
	session.Close()
	return
}
