package pr

import (
	"fmt"
	"log"
	"strings"
)

func Synced(c Client, org, repo, id string) ([]PullRequest, error) {

	// Get the url for the PR
	pr, err := getPr(c, org, repo, id)

	if err != nil {
		log.Fatal(err)
	}

	// Find PRs that mention the url
	repos := []repoFilter{
		buildRepoFilter(WordPressOrg, "gutenberg", "is:open", "is:pr", `label:"Mobile App - i.e. Android or iOS"`),
		buildRepoFilter(WpMobileOrg, "gutenberg-mobile", "is:open", "is:pr"),
		buildRepoFilter(WpMobileOrg, "WordPress-Android", "is:open", "is:pr", `label:"Gutenberg"`),
		buildRepoFilter(WpMobileOrg, "WordPress-iOS", "is:open", "is:pr", `label:"Gutenberg"`),
		buildRepoFilter(AutomatticOrg, "jetpack", "is:open", "is:pr"),
	}

	var foundPrs []PullRequest
	prChannel := make(chan SearchResult)

	orgRepo := fmt.Sprintf("%s/%s", org, repo)
	for _, r := range repos {
		// Ignore the repo we're searching
		if r.repo == orgRepo {
			continue
		}

		go func(r repoFilter) {
			result, err := searchPRs(c, r)

			if err != nil {
				fmt.Println(err)
			}
			prChannel <- result

		}(r)
	}

	for i := 1; i < len(repos); i++ {
		result := <-prChannel
		foundPrs = append(foundPrs, filterSyncedPrs(pr.Url, result.Items)...)
	}

	return foundPrs, nil
}

func filterSyncedPrs(target string, prs []PullRequest) []PullRequest {
	var filtered []PullRequest

	for _, pr := range prs {
		if strings.Contains(pr.Body, target) {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}
