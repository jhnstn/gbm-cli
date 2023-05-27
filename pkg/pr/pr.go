package pr

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Client interface {
	Get(string, interface{}) error
}

type PullRequest struct {
	Number int
	Url    string `json:"html_url"`
	Body   string
	Title  string
	Labels []struct{ Name string }
	State  string
	User   struct {
		Login string
	}
	Draft     bool
	Mergeable bool
	Org       string
}

type SearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []PullRequest
}

type repoFilter struct {
	repo  string
	query string
}

var (
	WpMobileOrg   string
	WordPressOrg  string
	AutomatticOrg string
)

func init() {
	if gbmWpMobileOrg, ok := os.LookupEnv("GBM_WPMOBILE_ORG"); !ok {
		WpMobileOrg = "wordpress-mobile"
	} else {
		WpMobileOrg = gbmWpMobileOrg
	}

	if gbmWordPressOrg, ok := os.LookupEnv("GBM_WORDPRESS_ORG"); !ok {
		WordPressOrg = "WordPress"
	} else {
		WordPressOrg = gbmWordPressOrg
	}

	if gbmAutomatticOrg, ok := os.LookupEnv("GBM_AUTOMATTIC_ORG"); !ok {
		AutomatticOrg = "automattic"
	} else {
		AutomatticOrg = gbmAutomatticOrg
	}
}

func GetPr(client Client, repo, id string) (PullRequest, error) {
	var org string

	switch repo {
	case "gutenberg":
		org = WordPressOrg
	case "jetpack":
		org = AutomatticOrg
	case "gutenberg-mobile":
		fallthrough
	case "WordPress-Android":
		fallthrough
	case "WordPress-iOS":
		org = WpMobileOrg
	}

	pr, err := getPr(client, org, repo, id)
	if err != nil {
		return PullRequest{}, err
	}

	pr.Org = org
	return pr, nil
}

func getPr(client Client, org, repo, id string) (PullRequest, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%s", org, repo, id)
	response := PullRequest{}
	err := client.Get(endpoint, &response)

	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

func buildRepoFilter(org, repo string, queries ...string) repoFilter {
	var encoded []string
	queries = append(queries, fmt.Sprintf("repo:%s/%s", org, repo))

	for _, q := range queries {
		encoded = append(encoded, url.QueryEscape(q))
	}

	return repoFilter{
		repo:  fmt.Sprintf("%s/%s", org, repo),
		query: strings.Join(encoded, "+"),
	}
}

func searchPRs(client Client, filter repoFilter) (SearchResult, error) {
	endpoint := fmt.Sprintf("search/issues?q=%s", filter.query)
	response := SearchResult{}
	err := client.Get(endpoint, &response)

	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}
