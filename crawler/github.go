package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GitHubResponse struct {
	TotalCount        int    `json:"total_count"`
	IncompleteResults bool   `json:"incomplete_results"`
	Items             []Item `json:"items"`
}

type Item struct {
	gorm.Model
	ID       int    `gorm:"primaryKey" json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	// Owner            Owner    `json:"owner"`
	HTMLURL          string  `json:"html_url"`
	Description      string  `json:"description"`
	Fork             bool    `json:"fork"`
	URL              string  `json:"url"`
	ForksURL         string  `json:"forks_url"`
	KeysURL          string  `json:"keys_url"`
	CollaboratorsURL string  `json:"collaborators_url"`
	TeamsURL         string  `json:"teams_url"`
	HooksURL         string  `json:"hooks_url"`
	IssueEventsURL   string  `json:"issue_events_url"`
	EventsURL        string  `json:"events_url"`
	AssigneesURL     string  `json:"assignees_url"`
	BranchesURL      string  `json:"branches_url"`
	TagsURL          string  `json:"tags_url"`
	BlobsURL         string  `json:"blobs_url"`
	GitTagsURL       string  `json:"git_tags_url"`
	GitRefsURL       string  `json:"git_refs_url"`
	TreesURL         string  `json:"trees_url"`
	StatusesURL      string  `json:"statuses_url"`
	LanguagesURL     string  `json:"languages_url"`
	StargazersURL    string  `json:"stargazers_url"`
	ContributorsURL  string  `json:"contributors_url"`
	SubscribersURL   string  `json:"subscribers_url"`
	SubscriptionURL  string  `json:"subscription_url"`
	CommitsURL       string  `json:"commits_url"`
	GitCommitsURL    string  `json:"git_commits_url"`
	CommentsURL      string  `json:"comments_url"`
	IssueCommentURL  string  `json:"issue_comment_url"`
	ContentsURL      string  `json:"contents_url"`
	CompareURL       string  `json:"compare_url"`
	MergesURL        string  `json:"merges_url"`
	ArchiveURL       string  `json:"archive_url"`
	DownloadsURL     string  `json:"downloads_url"`
	IssuesURL        string  `json:"issues_url"`
	PullsURL         string  `json:"pulls_url"`
	MilestonesURL    string  `json:"milestones_url"`
	NotificationsURL string  `json:"notifications_url"`
	LabelsURL        string  `json:"labels_url"`
	ReleasesURL      string  `json:"releases_url"`
	DeploymentsURL   string  `json:"deployments_url"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	PushedAt         string  `json:"pushed_at"`
	GitURL           string  `json:"git_url"`
	SSHURL           string  `json:"ssh_url"`
	CloneURL         string  `json:"clone_url"`
	SVNURL           string  `json:"svn_url"`
	Homepage         string  `json:"homepage"`
	Size             int     `json:"size"`
	StargazersCount  int     `json:"stargazers_count"`
	WatchersCount    int     `json:"watchers_count"`
	Language         string  `json:"language"`
	HasIssues        bool    `json:"has_issues"`
	HasProjects      bool    `json:"has_projects"`
	HasDownloads     bool    `json:"has_downloads"`
	HasWiki          bool    `json:"has_wiki"`
	HasPages         bool    `json:"has_pages"`
	HasDiscussions   bool    `json:"has_discussions"`
	ForksCount       int     `json:"forks_count"`
	MirrorURL        *string `json:"mirror_url"`
	Archived         bool    `json:"archived"`
	Disabled         bool    `json:"disabled"`
	OpenIssuesCount  int     `json:"open_issues_count"`
	// License          License `json:"license"`
	AllowForking     bool `json:"allow_forking"`
	IsTemplate       bool `json:"is_template"`
	WebCommitSignoff bool `json:"web_commit_signoff_required"`
	// Topics           []string `json:"topics"`
	Topics        pq.StringArray `gorm:"type:text[]" json:"topics"`
	Visibility    string         `json:"visibility"`
	Forks         int            `json:"forks"`
	OpenIssues    int            `json:"open_issues"`
	Watchers      int            `json:"watchers"`
	DefaultBranch string         `json:"default_branch"`
	Score         float64        `json:"score"`
}

type Owner struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

func getGitHubRepositories(db *gorm.DB, basePath string) {
	var allItems []Item
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

func getPerPage(page int, basePath string) (int, []Item) {
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
