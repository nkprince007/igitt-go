package github_test

import (
	"os"
	"testing"

	"github.com/ad2games/vcr-go"
	"github.com/nkprince007/igitt-go"
	"github.com/nkprince007/igitt-go/github"
	"github.com/stretchr/testify/assert"
)

func ok(got interface{}, want interface{}, method string, t *testing.T) {
	if !assert.Equal(t, got, want) {
		t.Errorf("%s was incorrect, got: '%v', want: '%v'", method, got, want)
	}
}

func TestRepository(t *testing.T) {
	var repo igitt.Repository

	vcr.Start("repo", nil)
	defer vcr.Stop()

	repo, err := github.NewRepository(1, os.Getenv("GITHUB_TEST_TOKEN"))
	if err != nil {
		t.Error(err)
	}

	if _, isOk := repo.(igitt.Repository); !isOk {
		t.Errorf("github.Repository does not implement all methods in the" +
			"interface igitt.Repository")
	}

	ok(repo.FullName(), "mojombo/grit", "FullName", t)
	ok(repo.Description(),
		"**Grit is no longer maintained. Check out libgit2/rugged.** Grit "+
			"gives you object oriented read/write access to Git repositories "+
			"via Ruby.",
		"Description", t)
	ok(repo.IsPrivate(), false, "IsPrivate", t)
	ok(repo.ID(), 1, "ID", t)
	ok(repo.WebURL(), "https://github.com/mojombo/grit", "WebURL", t)
	ok(repo.APIURL(), "https://api.github.com/repositories/1", "APIURL", t)
	ok(repo.Homepage(), "http://grit.rubyforge.org/", "Homepage", t)
	ok(repo.HasIssues(), true, "HasIssues", t)
	ok(repo.IsFork(), false, "IsFork", t)
}
