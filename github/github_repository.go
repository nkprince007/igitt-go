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
			return "https://api.github.com"
		}
		return url
	}()
}

type repository struct {
	ID          int         `json:"id"`
	FullName    string      `json:"full_name"`
	Description string      `json:"description"`
	Private     bool        `json:"private"`
	URL         string      `json:"url"`
	WebURL      string      `json:"html_url"`
	Homepage    string      `json:"homepage"`
	HasIssues   bool        `json:"has_issues"`
	IsFork      bool        `json:"fork"`
	Parent      *repository `json:"parent,omitempty"`
}

type label struct {
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

type labelArr struct {
	Labels []label
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

func (repo *Repository) decode() error {
	if repo.decoded {
		return nil
	}
	repo.decoded = true
	return json.Unmarshal(repo.data, &repo.r)
}

func (repo *Repository) refresh() error {
	if repo.fetched {
		repo.decode()
		return nil
	}

	url := repo.APIURL()
	data, e := make(chan []byte), make(chan error)
	go lib.Get(repo.token, url, data, e)
	if err := <-e; err != nil {
		return err
	}

	repo.data = <-data
	repo.decode()
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
	if repo.fetched {
		return repo.r.URL
	}
	url, _ := repo.url()
	return url
}

// FullName gives the fullname of GitHub repository
func (repo *Repository) FullName() string {
	if repo.fullName == "" {
		repo.refresh()
	}
	return repo.fullName
}

// ID returns the unique identifier of GitHub repository
func (repo *Repository) ID() int {
	if repo.id == 0 {
		repo.refresh()
	}
	return repo.id
}

// Description returns the description of the repository
func (repo *Repository) Description() string {
	repo.refresh()
	return repo.r.Description
}

// IsPrivate tells whether the GitHub repository is private or not
func (repo *Repository) IsPrivate() bool {
	repo.refresh()
	return repo.r.Private
}

// WebURL returns the URL for the webpage to GitHub repository
func (repo *Repository) WebURL() string {
	repo.refresh()
	return repo.r.WebURL
}

// Homepage returns the URL for homepage of GitHub repository
func (repo *Repository) Homepage() string {
	repo.refresh()
	return repo.r.Homepage
}

// HasIssues tells whether GitHub repository has an issue tracker or not
func (repo *Repository) HasIssues() bool {
	repo.refresh()
	return repo.r.HasIssues
}

// IsFork tells whether the GitHub repository is a fork or not
func (repo *Repository) IsFork() bool {
	repo.refresh()
	return repo.r.IsFork
}

// Parent returns the parent repository if it is a fork
func (repo *Repository) Parent() *Repository {
	if !repo.IsFork() {
		return nil
	}
	parent := repo.r.Parent
	return &Repository{
		parent, 0, parent.FullName, repo.token, nil, false, false}
}

// GetAllLabels returns all the labels associated with this repository
func (repo *Repository) GetAllLabels() ([]string, error) {
	labelArr := []label{}
	labels := []string{}

	data, e := make(chan []byte), make(chan error)
	go lib.Get(repo.token, repo.APIURL()+"/labels", data, e)
	if err := <-e; err != nil {
		return nil, err
	}

	if err := json.Unmarshal(<-data, &labelArr); err != nil {
		return nil, err
	}

	for _, l := range labelArr {
		labels = append(labels, l.Name)
	}
	return labels, nil
}

// CreateLabel creates a new label in this GitHub repository
func (repo *Repository) CreateLabel(
	name, color, description, labelType string) error {
	l := &label{name, color, description}
	data, e := make(chan []byte), make(chan error)

	go lib.Post(repo.token, repo.APIURL()+"/labels", data, e)
	encoded, err := json.Marshal(l)
	if err != nil {
		return err
	}
	data <- encoded
	err = <-e

	repo.decoded = false
	repo.data = <-data
	repo.decode()
	return err
}

// DeleteLabel deletes an existing label in the repository
func (repo *Repository) DeleteLabel(name string) error {
	data, e := make(chan []byte), make(chan error)
	go lib.Delete(repo.token, repo.APIURL()+"/labels/"+name, data, e)
	data <- nil
	return <-e
}

func (repo *Repository) String() string {
	return fmt.Sprintf("Repository(id=%d, fullName=%s, url=%s)",
		repo.id, repo.fullName, repo.APIURL())
}
