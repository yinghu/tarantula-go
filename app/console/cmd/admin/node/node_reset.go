package node

import (
	"fmt"

	"gameclustering.com/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	resetCmd.Flags().StringP("env", "E", "dev", "env")
	resetCmd.Flags().String("etcd", "192.168.1.7:2379", "etcd host")
	resetCmd.Flags().StringP("app", "A", "", "app (required)")
	resetCmd.MarkFlagRequired("app")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset node",
	Long:  "reset node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("etcd")
		app, _ := cmd.Flags().GetString("app")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		prefix := fmt.Sprintf("%s/node", env)
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			nidKey := fmt.Sprintf("app.%s.nid", app)
			err := ctx.Del(nidKey, false)
			if err != nil {
				return err
			}
			return ctx.Del(app, true)
		})
		if err != nil {
			fmt.Printf("reset command failed %s\n", err.Error())
		}
	},
}
