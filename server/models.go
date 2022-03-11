package main

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

type Bot struct {
	Country string
	Ip string
	OnlineAt time.Time
	Online bool
	gorm.Model
}

func ConnectBot(ip string){
	var bot Bot
	r := DB.First(&bot, "ip = ?", ip)
	if r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		log.Print(r.Error)
	}

	if r.RowsAffected == 0 {
		DB.Create(&Bot{
			Country: getCountry(ip),
			Ip: ip,
			Online: true,
			OnlineAt: time.Now(),
		})
	}
}


func SetupModels() {
	var err error
	DB, err = gorm.Open(sqlite.Open("bots.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(
		&Bot{},
	)

}
