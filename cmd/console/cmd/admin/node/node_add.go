package node

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().StringP("name", "N", "", "name (required)")
	addCmd.MarkFlagRequired("name")
	//addCmd.Flags().StringP("password", "P", "", "password (required)")
	//addCmd.MarkFlagRequired("password")
	//addCmd.Flags().StringP("host", "H", "http://localhost", "use localhost")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add node",
	Long:  "add node",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fmt.Printf("add node %s\n", name)
	},
}
