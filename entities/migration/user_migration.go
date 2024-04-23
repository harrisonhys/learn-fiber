package migration

import (
	"fmt"

	"github.com/harrisonhys/learn-fiber/config"
	"github.com/harrisonhys/learn-fiber/entities"
)

func MigrateUser() {
	err := config.DB.AutoMigrate(&entities.User{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrated user table")
}
