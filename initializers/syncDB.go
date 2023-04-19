package initializers

import "github.com/adamlahbib/gitaz/models"

func SyncDB() {
	DB.AutoMigrate(&models.Repo{}, &models.User{}, &models.Deployment{})
}
