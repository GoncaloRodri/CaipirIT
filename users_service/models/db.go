package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    var err error
    DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    err = DB.AutoMigrate(&User{}, &Organization{}, &Team{}, &Membership{}, &OrgAdmin{})
    if err != nil {
        log.Fatal("Migration failed:", err)
    }
}