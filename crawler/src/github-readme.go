package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getRepoReadme() {
	fmt.Println("Running getRepoReadme function")
	// Database connection details
	db := getStorage()

	// Pagination parameters
	page := 1
	pageSize := 10
	for {
		offset := (page - 1) * pageSize
		fmt.Printf("Reading repositories with pagination: page=%d, pageSize=%d, offset=%d\n", page, pageSize, offset)

		// Read repositories with pagination
		var repositories []Repository
		if err := db.Limit(pageSize).Offset(offset).Find(&repositories).Error; err != nil {
			log.Fatalf("failed to read repositories: %v", err)
		}

		if len(repositories) == 0 {
			break
		}

		run(db, repositories)
		page++
	}

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

func run(db *gorm.DB, repositories []Repository) {

	// Read all repositories from the repositories table
	// var repositories []Repository
	// if err := db.Find(&repositories).Error; err != nil {
	// 	log.Fatalf("failed to read repositories: %v", err)
	// }

	// Print the retrieved repositories
	for _, repo := range repositories {
		fmt.Printf("ID: %d, Name: %s, FullName: %s\n", repo.ID, repo.Name, repo.FullName)

		urls := getReadmeUrls(repo.FullName, repo.DefaultBranch)
		readme, err := fetchReadmeText(urls[0])
		if err != nil {
			log.Printf("\tfailed to fetch README for %s %s: %v", repo.FullName, urls[0], err)
			readme, err = fetchReadmeText(urls[1])
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Printf("\tfailed to fetch README for %s %s: %v", repo.FullName, urls[1], err)
			time.Sleep(2 * time.Second)
			continue
		}

		summary := RepositorySummary{
			RepositoryID:    repo.ID,
			Name:            repo.Name,
			FullName:        repo.FullName,
			Size:            repo.Size,
			StargazersCount: repo.StargazersCount,
			DefaultBranch:   repo.DefaultBranch,
			OpenIssuesCount: repo.OpenIssuesCount,
			Language:        repo.Language,
			PushedAt:        repo.PushedAt,
			ReadMe:          readme,
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
		time.Sleep(2 * time.Second)

		// https://raw.githubusercontent.com/myshell-ai/OpenVoice/refs/heads/main/README.md
	}
}

func getReadmeUrls(fullName string, defaultBranch string) []string {
	return []string{
		fmt.Sprintf("https://raw.githubusercontent.com/%s/refs/heads/%s/README.md", fullName, defaultBranch),
		fmt.Sprintf("https://raw.githubusercontent.com/%s/refs/heads/%s/README.markdown", fullName, defaultBranch),
	}
}

func fetchReadmeText(url string) (string, error) {
	fmt.Printf("\tFetching README from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("\tfailed to fetch README: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
