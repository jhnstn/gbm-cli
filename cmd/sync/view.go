package sync

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/wordpress-mobile/gbm/pkg/pr"
)

// releaseCmd represents the release command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View synced prs",
	Long: `
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.DefaultRESTClient()
		if err != nil {
			return err
		}
		var (
			repo string
			id   string
		)

		if len(args) == 1 {
			if i, ok := strconv.Atoi(args[0]); ok != nil {
				p := strings.Split(args[0], "/")
				if len(p) != 2 {
					return fmt.Errorf("invalid pr format")
				}
				repo = p[0]
				id = p[1]

			} else {
				id = fmt.Sprint(i)
			}
		}

		if len(args) == 2 {
			repo = args[0]
			id = args[1]
		}

		p, err := pr.GetPr(client, repo, id)

		if err != nil {
			return err
		}
		fmt.Printf("Checking %s\n", p.Url)

		synced, err := pr.Synced(client, p.Org, repo, id)
		if err != nil {
			return err
		}

		fmt.Printf("Found %d synced prs\n\n", len(synced))

		headerFmt := color.New(color.Bold, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgCyan).SprintfFunc()
		tbl := table.New("ID", "URL", "State", "Mergeable")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, p := range append([]pr.PullRequest{p}, synced...) {
			tbl.AddRow(p.Number, p.Url, p.State, p.Mergeable)
		}

		tbl.Print()

		return nil
	},
}

func init() {
}
