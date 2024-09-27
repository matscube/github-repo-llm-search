package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// todo: use os.Getenv to get the environment variables
	host := "0.0.0.0"
	user := "exampleuser"
	password := "examplepass"
	dbname := "exampledb"

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&Item{})
	if err != nil {
		log.Fatal(err)
	}

	// Fetch and store GitHub repositories
	getGitHubRepositories(db)

	// getGitHubRepositories()
}
