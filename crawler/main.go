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

	var basePaths []string
	// > 30000 stars
	// basePaths = append(basePaths, ">30000")
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
	basePaths = append(basePaths, getRangeWithSlidingWindow(500, 1000, 3, 2, true)...)

	fmt.Println(basePaths)

	for _, basePath := range basePaths {
		url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:%s&sort=stars", basePath)
		getGitHubRepositories(db, url)
	}
	// Fetch and store GitHub repositories

	// getGitHubRepositories()
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
