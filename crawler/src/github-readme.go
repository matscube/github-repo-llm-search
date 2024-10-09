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
		if err := db.Order("id ASC").Limit(pageSize).Offset(offset).Find(&repositories).Error; err != nil {
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

		// Check if a RepositorySummary record already exists for this repository
		var existingSummary RepositorySummary
		if err := db.Where("repository_id = ?", repo.ID).First(&existingSummary).Error; err == nil {
			fmt.Printf("\tRepositorySummary already exists for %s, skipping...\n", repo.FullName)
			continue
		}

		notFoundList := readmeNotFoundRepos()
		if contains(notFoundList, repo.FullName) {
			fmt.Printf("\tREADME not found for %s, skipping...\n", repo.FullName)
			continue
		}

		urls := getReadmeUrls(repo.FullName, repo.DefaultBranch)
		readme, err := fetchReadmeTextThroughUrls(urls)
		if err != nil {
			fmt.Printf("\tfailed to fetch README for %s\n", fmt.Sprintf("https://github.com/%s", repo.FullName))
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
			fmt.Printf("\tfailed to create new repository record for %s: %v\n", repo.FullName, err)
		} else {
			fmt.Printf("\tSuccessfully created new repository record for %s\n", repo.FullName)
		}
		time.Sleep(2 * time.Second)

		// https://raw.githubusercontent.com/myshell-ai/OpenVoice/refs/heads/main/README.md
	}
}

func fetchReadmeTextThroughUrls(urls []string) (string, error) {
	for _, url := range urls {
		readme, err := fetchReadmeText(url)
		if err == nil {
			return readme, nil
		}
		fmt.Printf("\t\tfailed to fetch README for %s: %v\n", url, err)
		time.Sleep(1 * time.Second)
	}
	return "", fmt.Errorf("\t\tfailed to fetch README from all URLs\n")
}

func readmeNotFoundRepos() []string {
	return []string{
		"mig/gedit-themes",
		"blynn/gitmagic",
		"planetbeing/iphonelinux",
		"luabind/luabind",
	}
}

func getReadmeUrls(fullName string, defaultBranch string) []string {
	baseUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/refs/heads/%s/", fullName, defaultBranch)
	files := []string{
		"README.md",
		"readme.md",
		"Readme.md",
		"README.rst",
		"README",
		"readme",
		"README.rdoc",
		"README.textile",
		"README.markdown",
		"Readme.markdown",
		"README.mkd",
		"README.mkdn",
		"readme.html",
	}
	var urls []string
	for _, file := range files {
		urls = append(urls, baseUrl+file)
	}
	return urls
}

func fetchReadmeText(url string) (string, error) {
	fmt.Printf("\tFetching README from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("\tfailed to fetch README: status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
