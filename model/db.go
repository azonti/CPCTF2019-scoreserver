package model

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
)

var mgoSess *mgo.Session
var db *mgo.Database

//InitDB Initialize Database
func InitDB() error {
	sess, err := mgo.Dial(os.Getenv("MONGODB_URL"))
	if err != nil {
		return fmt.Errorf("failed to establish DB session: %v", err)
	}
	mgoSess = sess
	db = mgoSess.DB("")
	return nil
}

//TermDB Terminate Database
func TermDB() {
	mgoSess.Close()
	return
}
