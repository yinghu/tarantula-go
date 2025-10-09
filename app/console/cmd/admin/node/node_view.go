package node

import (
	"fmt"

	"gameclustering.com/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	viewCmd.Flags().StringP("env", "E", "dev", "env")
	viewCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
	viewCmd.Flags().StringP("app", "A", "", "app (required)")
	viewCmd.MarkFlagRequired("app")
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view node",
	Long:  "view node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		app, _ := cmd.Flags().GetString("app")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		prefix := fmt.Sprintf("%s/node", env)
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			err := ctx.List(app, func(k, v string) bool{
				fmt.Printf("%s : %s\n", k, v)
				return true
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("view command failed %s\n", err.Error())
		}
	},
}
