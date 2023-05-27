package sync

import (
	"fmt"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var RootCmd = &cobra.Command{
	Use:   "sync",
	Short: "TBD",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync called")
	},
}

func init() {

	RootCmd.AddCommand(viewCmd)
}
