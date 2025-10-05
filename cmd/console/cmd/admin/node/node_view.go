package node

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	viewCmd.Flags().StringP("name", "N", "", "name (required)")
	viewCmd.MarkFlagRequired("name")
	//addCmd.Flags().StringP("password", "P", "", "password (required)")
	//addCmd.MarkFlagRequired("password")
	//addCmd.Flags().StringP("host", "H", "http://localhost", "use localhost")
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view node",
	Long:  "view node",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fmt.Printf("view node %s\n", name)
	},
}
