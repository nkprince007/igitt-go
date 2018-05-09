package igitt

import "github.com/nkprince007/igitt-go/github"

// Repository represents a repository on any Git hosting provider like GitHub,
// GitLab, etc.
type Repository interface {
	FullName() string
	ID() int
	Description() string
	WebURL() string
	APIURL() string
	Homepage() string
	HasIssues() bool
	IsPrivate() bool
	IsFork() bool
	Parent() *github.Repository

	// Delete()
	// Clone()
	// Fork(org string, namespace string)
	// GetPermissionLevel(username string)

	// RegisterWebhook(url string, secret string, events []string)
	// UnregisterWebhook(url string)
	// Hooks()

	// CreateIssue(title string, body string)
	// GetIssue(number int)
	// FilterIssues(
	// 	createdAfter, createdBefore, updatedAfter, updatedBefore time.Time,
	// 	state string)
	// GetAllIssues()

	// CreateMergeRequest(
	// 	title, base, head, body, targetProject string, targetProjectID int)
	// GetMergeRequest(number int)
	// FilterMergeRequests(
	// 	createdAfter, createdBefore, updatedAfter, updatedBefore time.Time,
	// 	state string)
	// GetAllMergeRequests()

	CreateLabel(
		name string, color string, description string, labelType string) error
	DeleteLabel(name string) error
	GetAllLabels() ([]string, error)

	// GetAllCommits()
	String() string
}

// type repository struct {
// 	Identifier  int
// 	FullName    string
// 	CloneURL    string
// 	Parent      *repository
// 	TopLevelOrg string
// }
