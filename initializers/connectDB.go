package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := os.Getenv("DB_DSN")

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		panic("failed to connect database")
	}

}
