package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().IntP("port", "P", 8090, "use 8080")
	versionCmd.Flags().StringP("host", "H", "localhost", "use localhost")
}

var versionCmd = &cobra.Command{
	Use:   "dial",
	Short: "dial tarantula",
	Long:  "dial tarantula over tcp",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		host, _ := cmd.Flags().GetString("host")
		fmt.Printf("dialing %s:%d\n", host, port)
	},
}
