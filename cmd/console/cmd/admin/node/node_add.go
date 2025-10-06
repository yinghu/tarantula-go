package node

import (
	"fmt"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().StringP("env", "E", "dev", "env")
	addCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
	addCmd.Flags().StringP("app", "A", "", "app (required)")
	addCmd.MarkFlagRequired("app")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add node",
	Long:  "add node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		app, _ := cmd.Flags().GetString("app")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		prefix := fmt.Sprintf("%s.%s", env,app)
		cnf := conf.Config{Sequence: 2}
		cx.Execute(prefix,func(ctx core.Ctx) error {
			name := fmt.Sprintf("%s.%d",app,cnf.Sequence)
			return ctx.Put(name, string(util.ToJson(cnf)))
		})
		fmt.Printf("add node %s %s %s\n", prefix, env, host)
	},
}
