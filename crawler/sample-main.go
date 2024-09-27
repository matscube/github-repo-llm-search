package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

const (
	numWorkers = 10 // Number of concurrent workers
)

func sampleMain() {
	// Set up GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Search for repositories
	repos, _, err := client.Search.Repositories(ctx, "stars:>10000", &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		log.Fatalf("Error searching repositories: %v", err)
	}

	// Create a channel to send repositories to workers
	repoChan := make(chan *github.Repository, len(repos.Repositories))

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, client, repoChan, &wg)
	}

	// Send repositories to the channel
	for _, repo := range repos.Repositories {
		repoChan <- repo
	}
	close(repoChan)

	// Wait for all workers to finish
	wg.Wait()
}

func worker(ctx context.Context, client *github.Client, repoChan <-chan *github.Repository, wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range repoChan {
		fmt.Printf("Processing repository: %s\n", *repo.FullName)

		// Get repository metadata
		repoMeta, _, err := client.Repositories.Get(ctx, *repo.Owner.Login, *repo.Name)
		if err != nil {
			log.Printf("Error getting metadata for %s: %v", *repo.FullName, err)
			continue
		}

		fmt.Printf("Repository metadata for %s:\n%+v\n", *repo.FullName, repoMeta)
	}
}
