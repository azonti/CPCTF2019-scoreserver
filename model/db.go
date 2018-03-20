package model

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
)

var mgoSession *mgo.Session
var db *mgo.Database

//InitDB Initialize Database
func InitDB() error {
	var err error
	mgoSession, err = mgo.Dial(os.Getenv("MONGODB_URL"))
	if err != nil {
		return fmt.Errorf("failed to establish DB session")
	}
	db = mgoSession.DB("")
	return nil
}

//TermDB Terminate Database
func TermDB() {
	mgoSession.Close()
	return
}
