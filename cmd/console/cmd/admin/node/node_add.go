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
		prefix, _ := cmd.Flags().GetString("app")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		//lkey := fmt.Sprintf("%s.config", env)
		cx.Execute("config", func(ctx core.Ctx) error {
			//idx := ctx.AppIndex(env)
			//slices.SortFunc(idx.Index,func(a,b) int {})
			cnf := conf.Config{Sequence: 1}
			//fmt.Printf("atomic on %s %d\n", lkey, cnf.Sequence)
			v, err := ctx.Get(prefix)
			if err == nil {
				fmt.Printf("already existed %s\n", v)
				return err
			}
			return ctx.Put(prefix, string(util.ToJson(cnf)))
		})
		fmt.Printf("add node %s %s %s\n", prefix, env, host)
	},
}
