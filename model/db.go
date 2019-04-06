package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
	"os"
)

var db *gorm.DB

//InitDB Initialize Database
func InitDB() error {
	var err error
	db, err = gorm.Open("mysql", os.Getenv("MARIADB_URL")+"?parseTime=True")
	if err != nil {
		return fmt.Errorf("failed to connect DB: %v", err)
	}
	db = db.Set("gorm:save_associations", false)
	if err := db.AutoMigrate(&Challenge{}, &Hint{}, &Vote{}, &Question{}, &User{}).Error; err != nil {
		return err
	}
	return nil
}

//TermDB Terminate Database
func TermDB() {
	db.Close()
	return
}
