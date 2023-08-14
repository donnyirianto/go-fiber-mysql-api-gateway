package models

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load configuration from environment variables or config file
	loadConfig()

	dsn := viper.GetString("DB_DSN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func loadConfig() {
	viper.SetConfigFile(".env") //path ke file .env
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	viper.AutomaticEnv()

	viper.SetDefault("DB_DSN", "root:n3wbi329m3d@tcp(127.0.0.1:3306)/edpreg4?charset=utf8mb4&parseTime=True&loc=Local")
}
