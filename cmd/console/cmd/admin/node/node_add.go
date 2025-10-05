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
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		name, _ := cmd.Flags().GetString("name")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		lkey := fmt.Sprintf("%s.config", env)
		cx.Execute(lkey, func(ctx core.Ctx) error {
			cnf := conf.Config{Sequence: 1}
			fmt.Printf("atomic on %s %d\n", lkey, cnf.Sequence)
			v, err := ctx.Get(name)
			if err == nil {
				fmt.Printf("already existed %s\n", v)
				return err
			}
			return ctx.Put(name, string(util.ToJson(cnf)))
		})
		fmt.Printf("add node %s %s %s\n", name, env, host)
	},
}
