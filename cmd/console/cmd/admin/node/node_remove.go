package node

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.Flags().StringP("name", "N", "", "name (required)")
	removeCmd.MarkFlagRequired("name")
	//addCmd.Flags().StringP("password", "P", "", "password (required)")
	//addCmd.MarkFlagRequired("password")
	//addCmd.Flags().StringP("host", "H", "http://localhost", "use localhost")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove node",
	Long:  "remove node",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fmt.Printf("remove node %s\n", name)
	},
}
