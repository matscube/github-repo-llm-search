package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define subcommands
	crawlCmd := flag.NewFlagSet("crawl", flag.ExitOnError)
	repoCmd := crawlCmd.Bool("repo", false, "Run getGitHubRepositories function")
	readmeCmd := crawlCmd.Bool("readme", false, "Run getRepoReadme function")

	// Check if at least one subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("expected 'crawl' subcommand")
		os.Exit(1)
	}

	// Parse the subcommand
	switch os.Args[1] {
	case "crawl":
		crawlCmd.Parse(os.Args[2:])
	default:
		fmt.Println("expected 'crawl' subcommand")
		os.Exit(1)
	}

	// Handle the subcommands
	if crawlCmd.Parsed() {
		if *repoCmd {
			getGitHubRepositories()
		} else if *readmeCmd {
			getRepoReadme()
		} else {
			fmt.Println("expected 'repo' or 'readme' flag")
			os.Exit(1)
		}
	}
}

func getRepoReadme() {
	fmt.Println("Running getRepoReadme function")
}
