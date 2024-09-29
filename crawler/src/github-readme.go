package main

import (
	"fmt"
	"log"

	"gorm.io/gorm/clause"
)

func getRepoReadme() {
	fmt.Println("Running getRepoReadme function")
	run()
	// db := getStorage()

	// summary := RepositorySummary{
	// 	Name:         "example-summary",
	// 	FullName:     "exampleuser/example-repo-summary",
	// 	RepositoryID: 26}
	// // english: I think that you need to specify columns for conflicts other than the primary key.
	// db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "repository_id"}}, // Specify the conflict target
	// 	UpdateAll: true,
	// }).Create(&summary)

}

func run() {
	// Database connection details
	db := getStorage()

	// Read all repositories from the repositories table
	// var repositories []Repository
	// if err := db.Find(&repositories).Error; err != nil {
	// 	log.Fatalf("failed to read repositories: %v", err)
	// }

	// Pagination parameters
	page := 1
	pageSize := 10
	offset := (page - 1) * pageSize

	// Read repositories with pagination
	var repositories []Repository
	if err := db.Limit(pageSize).Offset(offset).Find(&repositories).Error; err != nil {
		log.Fatalf("failed to read repositories: %v", err)
	}

	// Print the retrieved repositories
	for _, repo := range repositories {
		fmt.Printf("ID: %d, Name: %s, FullName: %s\n", repo.ID, repo.Name, repo.FullName)
		summary := RepositorySummary{
			RepositoryID: repo.ID,
			Name:         repo.Name,
			FullName:     repo.FullName,
		}
		// english: I think that you need to specify columns for conflicts other than the primary key.
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "repository_id"}}, // Specify the conflict target
			UpdateAll: true,
		}).Create(&summary).Error; err != nil {
			log.Printf("\tfailed to create new repository record for %s: %v", repo.FullName, err)
		} else {
			fmt.Printf("\tSuccessfully created new repository record for %s\n", repo.FullName)
		}

	}
}
