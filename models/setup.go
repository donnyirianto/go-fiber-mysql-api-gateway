package models

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load configuration from environment variables or config file
	loadConfig()

	dsn := viper.GetString("DB_DSN")

	// Create a new logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Custom  log output
		logger.Config{
			SlowThreshold:             200, // Set the threshold untuk slow SQL queries ( milliseconds)
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, // Buat custom logger
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	// Auto-migrate
	err = DB.AutoMigrate(&Log{})
	if err != nil {
		log.Fatal(err)
	}
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
