package cmd

import (
	"gameclustering.com/cmd/player"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tarantula",
		Short: "tarantula",
		Long:  "tarantula",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(player.PlayerCmd)
}
