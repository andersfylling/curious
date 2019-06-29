package curious

import (
	"context"
	"errors"
	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"strings"
)

type ProjectInfo struct {
	Name   string
	Author string
	URL    url.URL
	Stars  uint
}

func (p *ProjectInfo) String() string {
	return p.URL.String()
}

func GithubSearch(depName string) (projects []*ProjectInfo, err error) {
	splits := strings.Split(depName, "/")
	if len(splits) < 3 {
		return nil, errors.New("must have at least two slashes in a git url to id owner and project name")
	}
	//source := splits[0]
	owner := splits[1]
	project := splits[2]

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("CURIOUS_GITHUB_TOKEN")},
	)
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, src)
	client := github.NewClient(httpClient)

	var results []*github.CodeResult
	options := github.SearchOptions{}
	options.PerPage = 100
	for {
		result, _, err := client.Search.Code(ctx, depName, &options)
		if err != nil {
			return nil, errors.New("when executing GET request - " + err.Error())
		}

		if len(result.CodeResults) < options.PerPage {
			break
		}
		options.Page++

		for i := range result.CodeResults {
			if !strings.HasSuffix(*result.CodeResults[i].Path, ".go") {
				continue
			}
			results = append(results, &result.CodeResults[i])
		}
	}

	registerredProject := func(u *url.URL) bool {
		for i := range projects {
			if projects[i].URL.String() == u.String() {
				return true
			}
		}

		return false
	}

	for i := range results {
		r := results[i].Repository
		usr := r.GetOwner()
		u, err := url.Parse(r.GetHTMLURL())
		if err != nil {
			continue
		}
		if registerredProject(u) {
			continue
		}
		if usr != nil && usr.GetName() != owner && r.GetName() != project {
			projects = append(projects, &ProjectInfo{
				Name:   r.GetName(),
				Author: r.GetOwner().GetName(),
				URL:    *u,
				Stars:  uint(r.GetStargazersCount()),
			})
		}
	}

	return projects, nil
}
