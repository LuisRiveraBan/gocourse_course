package bootstrap

import (
	"fmt"
	"github.com/LuisRiveraBan/gocourse_domain/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// InitLogger initializes a new logger for the application.
// It logs to stdout with timestamps and the file name and line number.
//
// Returns a pointer to the initialized logger.
//
// Example:
func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func ConnectToDatabase() (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	// Initialize the database connection and migrate the schema
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	// Enable logging for debugging purposes.
	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	// Migrate the schema for the User model in the database.
	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.Course{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}
