package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/nkprince007/igitt-go/lib"
)

var gitHubBaseURL string

func init() {
	gitHubBaseURL = func() string {
		url, exists := os.LookupEnv("GITHUB_BASE_URL")
		if !exists {
			gitHubBaseURL = "https://api.github.com"
		}
		return url
	}()
}

type repository struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	URL         string `json:"url"`
	WebURL      string `json:"html_url"`
	Homepage    string `json:"homepage"`
	HasIssues   bool   `json:"has_issues"`
	IsFork      bool   `json:"fork"`
}

// Repository represents a repository on GitHub
type Repository struct {
	r        *repository
	id       int
	fullName string
	token    string
	data     []byte
	fetched  bool
	decoded  bool
}

func (repo *Repository) url() (string, error) {
	if repo.id != 0 {
		return fmt.Sprintf("%s/repositories/%d", gitHubBaseURL, repo.id), nil
	}
	if repo.fullName != "" {
		return fmt.Sprintf("%s/repos/%s", gitHubBaseURL, repo.fullName), nil
	}
	return "", errors.New("URL couldn't be formed as neither Identifier nor" +
		" FullName were specified")
}

func (repo *Repository) decodeFromData() error {
	if repo.decoded {
		return nil
	}
	repo.decoded = true
	return json.Unmarshal(repo.data, &repo.r)
}

func (repo *Repository) fetch() error {
	if repo.fetched {
		repo.decodeFromData()
		return nil
	}
	url, err := repo.url()
	if err != nil {
		return err
	}
	resp, err := lib.Get(repo.token, url)
	if err != nil {
		return err
	}
	repo.data = resp
	repo.decodeFromData()
	repo.id, repo.fullName, repo.fetched = repo.r.ID, repo.r.FullName, true
	return nil
}

// NewRepository returns a new GitHub Repository object
func NewRepository(id int, token string) (*Repository, error) {
	if id == 0 {
		return nil, errors.New("Repository: ID cannot be zero")
	}
	return &Repository{&repository{}, id, "", token, nil, false, false}, nil
}

// NewRepositoryFromName returns a new GitHub Repository object
func NewRepositoryFromName(name, token string) (*Repository, error) {
	if name == "" {
		return nil, errors.New("Repository: FullName cannot be an empty string")
	}
	return &Repository{&repository{}, 0, name, token, nil, false, false}, nil
}

// APIURL returns the API URL for the GitHub repository
func (repo *Repository) APIURL() string {
	url, err := repo.url()
	if err != nil && !repo.fetched {
		repo.fetch()
		return repo.r.URL
	}
	return url
}

// FullName gives the fullname of GitHub repository
func (repo *Repository) FullName() string {
	if repo.fullName == "" {
		repo.fetch()
	}
	return repo.fullName
}

// ID returns the unique identifier of GitHub repository
func (repo *Repository) ID() int {
	if repo.id == 0 {
		repo.fetch()
	}
	return repo.id
}

// Description returns the description of the repository
func (repo *Repository) Description() string {
	repo.fetch()
	return repo.r.Description
}

// IsPrivate tells whether the GitHub repository is private or not
func (repo *Repository) IsPrivate() bool {
	repo.fetch()
	return repo.r.Private
}

// WebURL returns the URL for the webpage to GitHub repository
func (repo *Repository) WebURL() string {
	repo.fetch()
	return repo.r.WebURL
}

// Homepage returns the URL for homepage of GitHub repository
func (repo *Repository) Homepage() string {
	repo.fetch()
	return repo.r.Homepage
}

// HasIssues tells whether GitHub repository has an issue tracker or not
func (repo *Repository) HasIssues() bool {
	repo.fetch()
	return repo.r.HasIssues
}

// IsFork tells whether the GitHub repository is a fork or not
func (repo *Repository) IsFork() bool {
	repo.fetch()
	return repo.r.IsFork
}

func (repo *Repository) String() string {
	return fmt.Sprintf("Repository(id=%d, fullName=%s, data=%v)",
		repo.id, repo.fullName, repo.r)
}
