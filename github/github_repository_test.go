package github_test

import (
	"os"
	"sort"
	"testing"

	"github.com/ad2games/vcr-go"
	"github.com/nkprince007/igitt-go"
	"github.com/nkprince007/igitt-go/github"
	"github.com/stretchr/testify/assert"
)

func check(got interface{}, want interface{}, method string, t *testing.T) {
	if !assert.Equal(t, want, got) {
		t.Errorf("%s was incorrect, got: '%v', want: '%v'", method, got, want)
	}
}

func TestRepository(t *testing.T) {
	var repo igitt.Repository

	vcr.Start("repo", nil)
	defer vcr.Stop()

	repo, err := github.NewRepositoryFromName(
		"nkprince007/coala-bears", os.Getenv("GITHUB_TEST_TOKEN"))
	if err != nil {
		t.Error(err)
	}

	if _, isOk := repo.(igitt.Repository); !isOk {
		t.Errorf("github.Repository does not implement all methods in the" +
			"interface igitt.Repository")
	}

	check(repo.String(),
		"Repository(id=0, fullName=nkprince007/coala-bears, "+
			"url=https://api.github.com/repos/nkprince007/coala-bears)",
		"String", t)
	check(repo.FullName(), "nkprince007/coala-bears", "FullName", t)
	check(repo.Description(), "Bears for coala", "Description", t)
	check(repo.IsPrivate(), false, "IsPrivate", t)
	check(repo.ID(), 76145200, "ID", t)
	check(repo.WebURL(), "https://github.com/nkprince007/coala-bears", "WebURL", t)
	check(repo.APIURL(),
		"https://api.github.com/repos/nkprince007/coala-bears", "APIURL", t)
	check(repo.Homepage(), "https://coala.io/", "Homepage", t)
	check(repo.HasIssues(), false, "HasIssues", t)
	check(repo.IsFork(), true, "IsFork", t)
	check(repo.Parent().FullName(), "coala/coala-bears", "Parent", t)

}

func TestRepositoryLabels(t *testing.T) {
	var repo igitt.Repository

	vcr.Start("repo-label", nil)
	defer vcr.Stop()

	repo, err := github.NewRepositoryFromName(
		"nkprince007/coala-bears", os.Getenv("GITHUB_TEST_TOKEN"))
	if err != nil {
		t.Error(err)
	}

	// check initial labels
	expectedLabels := []string{
		"bug", "duplicate", "enhancement", "help wanted", "invalid", "question",
		"test", "wontfix"}
	labels, err := repo.GetAllLabels()
	if err != nil {
		t.Error(err)
	}
	check(labels, expectedLabels, "GetAllLabels", t)

	// create new label
	err = repo.CreateLabel("test-again", "000000", "", "")
	if err != nil {
		t.Error(err)
	}

	// compare new label added
	labels, err = repo.GetAllLabels()
	if err != nil {
		t.Error(err)
	}
	newLabels := append(expectedLabels, "test-again")
	sort.Strings(newLabels)
	check(labels, newLabels, "CreateLabel", t)

	// delete new label
	err = repo.DeleteLabel("test-again")
	if err != nil {
		t.Error(err)
	}

	// compare old list of labels
	labels, err = repo.GetAllLabels()
	if err != nil {
		t.Error(err)
	}
	check(labels, expectedLabels, "DeleteLabel", t)
}
