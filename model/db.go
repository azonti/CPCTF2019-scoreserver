package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
	"os"
)

var db *gorm.DB

//InitDB Initialize the Database
func InitDB() error {
	var err error
	db, err = gorm.Open("mysql", os.Getenv("MARIADB_URL")+"?parseTime=True")
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&Challenge{}, &Hint{}, &Flag{}, &Vote{}, &Question{}, &User{}).Error; err != nil {
		return err
	}
	db = db.Set("gorm:save_associations", false)
	return nil
}

//TermDB Terminate the Database
func TermDB() {
	db.Close()
	return
}
