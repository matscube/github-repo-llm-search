package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GitHubResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []Repository `json:"items"`
}

func getGitHubRepositories() {
	db := getStorage()

	var basePaths []string
	// > 30000 stars
	basePaths = append(basePaths, ">30000")
	// > 20000 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(20000, 30000, 10000, 1000, true)...)
	// > 10000 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(10000, 20000, 1000, 100, true)...)
	// > 6000 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(6000, 10000, 200, 50, true)...)
	// > 2000 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(2000, 6000, 40, 10, true)...)
	// > 1000 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(1000, 2000, 10, 5, true)...)
	// > 500 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(500, 1000, 3, 2, true)...)
	// > 300 stars
	// basePaths = append(basePaths, getRangeWithSlidingWindow(300, 500, 1, 0, true)...)

	fmt.Println(basePaths)

	for _, basePath := range basePaths {
		url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:%s&sort=stars", basePath)
		getGitHubRepository(db, url)
	}
}

func getGitHubRepository(db *gorm.DB, basePath string) {
	var allItems []Repository
	page := 1
	failed := 0
	for {
		totalCount, items := getPerPage(page, basePath)
		fmt.Printf("totalCount: %d page: %d count: %d\n", totalCount, page, len(items))
		if len(items) == 0 {
			failed++
			fmt.Println("failed count: ", failed)
		}
		if failed > 3 {
			break
		}
		for i, item := range items {
			fmt.Printf("\t %d Item url: %s star: %d\n", len(allItems)+i+1, item.HTMLURL, item.StargazersCount)
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&item)
		}
		allItems = append(allItems, items...)
		page++
		time.Sleep(10 * time.Second) // Wait for 10 seconds before the next page
	}
	fmt.Println("Total count: ", len(allItems))
}

func getPerPage(page int, basePath string) (int, []Repository) {
	// Perform HTTP GET request
	// starThreshold := 1000
	perPage := 100 // max count 100
	url := fmt.Sprintf("%s&per_page=%d&page=%d", basePath, perPage, page)
	fmt.Println("url: ", url)

	// rate limit 10 req / min
	// api doc: https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28
	// https://api.github.com/search/repositories?q=stars:%3E1000&created>2020-07-01&sort=updated
	// url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:>%d&sort=stars&per_page=%d&page=%d", starThreshold, perPage, page)
	// url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:>%d&sort=stars&per_page=%d&page=%d", starThreshold, perPage, page)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error performing GET request: %v", err)
	}
	defer resp.Body.Close()

	var body GitHubResponse
	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	err = json.Unmarshal(responseData, &body)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return 0, nil
	}
	return body.TotalCount, body.Items
}
func getRangeWithSlidingWindow(start, end, perPage, overlap int, reversed bool) []string {
	var basePaths []string
	for {
		basePaths = append(basePaths, fmt.Sprintf("%d..%d", start, start+perPage+overlap))
		start += perPage
		if start >= end {
			break
		}
	}
	if reversed {
		return reverseStrings(basePaths)
	} else {
		return basePaths
	}
}

func reverseStrings(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
