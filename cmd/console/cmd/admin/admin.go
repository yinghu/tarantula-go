package admin

import (
	"github.com/spf13/cobra"
)

var (
	AdminCmd = &cobra.Command{
		Use:   "admin",
		Short: "tarantula admin",
		Long:  "tarantula admin",
	}
)

func Execute() error {
	return AdminCmd.Execute()
}

func init() {
	AdminCmd.AddCommand(loginCmd)
}
