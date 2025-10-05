package node

import (
	"github.com/spf13/cobra"
)

var (
	NodeCmd = &cobra.Command{
		Use:   "node",
		Short: "tarantula admin node",
		Long:  "tarantula admin node",
	}
)

func Execute() error {
	return NodeCmd.Execute()
}

func init() {
	NodeCmd.AddCommand(addCmd)
	NodeCmd.AddCommand(removeCmd)
	NodeCmd.AddCommand(viewCmd)
}
