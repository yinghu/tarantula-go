package player

import (
	"github.com/spf13/cobra"
)

var (
	PlayerCmd = &cobra.Command{
		Use:   "player",
		Short: "tarantula player",
		Long:  "tarantula player",
	}
)

func Execute() error {
	return PlayerCmd.Execute()
}

func init() {
	PlayerCmd.AddCommand(loginCmd)
	PlayerCmd.AddCommand(simuCmd)
}
