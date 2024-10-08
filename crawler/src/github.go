package main

import (
	"fmt"
	"log"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepositorySummary struct {
	gorm.Model
	ID              int        `gorm:"primaryKey"`
	Name            string     `gorm:"size:255"`
	FullName        string     `gorm:"size:255"`
	RepositoryID    int        `gorm:"not null; unique"`                      // Foreign key field
	Repository      Repository `gorm:"foreignKey:RepositoryID;references:ID"` // Foreign key relationship
	Size            int
	StargazersCount int
	DefaultBranch   string
	OpenIssuesCount int
	Language        string
	PushedAt        string
	ReadMe          string `gorm:"type:text"` // Use text data type for large text

}

type Repository struct {
	gorm.Model
	ID       int    `gorm:"primaryKey" json:"id"` // repo item id managed by github
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `gorm:"size:255;unique" json:"full_name"`
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

// type Owner struct {
// 	Login             string `json:"login"`
// 	ID                int    `json:"id"`
// 	NodeID            string `json:"node_id"`
// 	AvatarURL         string `json:"avatar_url"`
// 	GravatarID        string `json:"gravatar_id"`
// 	URL               string `json:"url"`
// 	HTMLURL           string `json:"html_url"`
// 	FollowersURL      string `json:"followers_url"`
// 	FollowingURL      string `json:"following_url"`
// 	GistsURL          string `json:"gists_url"`
// 	StarredURL        string `json:"starred_url"`
// 	SubscriptionsURL  string `json:"subscriptions_url"`
// 	OrganizationsURL  string `json:"organizations_url"`
// 	ReposURL          string `json:"repos_url"`
// 	EventsURL         string `json:"events_url"`
// 	ReceivedEventsURL string `json:"received_events_url"`
// 	Type              string `json:"type"`
// 	SiteAdmin         bool   `json:"site_admin"`
// }

// type License struct {
// 	Key    string `json:"key"`
// 	Name   string `json:"name"`
// 	SpdxID string `json:"spdx_id"`
// 	URL    string `json:"url"`
// 	NodeID string `json:"node_id"`
// }

func getStorage() *gorm.DB {
	// todo: use os.Getenv to get the environment variables
	host := "0.0.0.0"
	user := "exampleuser"
	password := "examplepass"
	dbname := "exampledb"

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&Repository{}, &RepositorySummary{})
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	// Example: Create a repository and a corresponding summary
	// repo := Repository{Name: "example-repo", FullName: "exampleuser/example-repo"}
	// db.Create(&repo)

	// summary := RepositorySummary{Name: "example-summary", FullName: "exampleuser/example-repo-summary", RepositoryID: repo.ID}
	// db.Create(&summary)

	return db
}
