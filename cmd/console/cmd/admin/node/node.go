package node

import (
	"github.com/spf13/cobra"
)

var (
	NodeCmd = &cobra.Command{
		Use:   "node",
		Short: "tarantula node",
		Long:  "tarantula node",
	}
)

func Execute() error {
	return NodeCmd.Execute()
}

func init() {
	NodeCmd.AddCommand(addCmd)
	NodeCmd.AddCommand(removeCmd)
}
