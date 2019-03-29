package model

import (
	"fmt"
	"os"
	"time"

	"github.com/globalsign/mgo"
)

var mgoSess *mgo.Session
var db *mgo.Database

//InitDB Initialize Database
func InitDB() error {
	mongoInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("MONGODB_URL")},
		Timeout:  20 * time.Second,
		Database: os.Getenv("MONGODB_DATABASE"),
		Username: os.Getenv("MONGODB_USERNAME"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	}
	sess, err := mgo.DialWithInfo(mongoInfo)
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
