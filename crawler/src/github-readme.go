package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getRepoReadme() {
	fmt.Println("Running getRepoReadme function")
	db := getStorage()

	summary := RepositorySummary{
		Name:         "example-summary",
		FullName:     "exampleuser/example-repo-summary",
		RepositoryID: 26}
	// english: I think that you need to specify columns for conflicts other than the primary key.
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "repository_id"}}, // Specify the conflict target
		UpdateAll: true,
	}).Create(&summary)

}

func run() {
	// Database connection details
	dsn := "host=localhost user=exampleuser dbname=exampledb password=examplepass sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Read all repositories from the repositories table
	var repositories []Repository
	if err := db.Find(&repositories).Error; err != nil {
		log.Fatalf("failed to read repositories: %v", err)
	}

	// Loop through each repository and create a new record in the new_repositories table
	// for _, repo := range repositories {
	// 	newRepo := RepositorySummary{
	// 		ID:       repo.ID,
	// 		Name:     repo.Name,
	// 		FullName: repo.FullName,
	// 	}
	// 	if err := db.Create(&newRepo).Error; err != nil {
	// 		log.Printf("failed to create new repository record for %s: %v", repo.FullName, err)
	// 	} else {
	// 		fmt.Printf("Successfully created new repository record for %s\n", repo.FullName)
	// 	}
	// }
}
