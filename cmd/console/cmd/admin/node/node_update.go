package node

import (
	"encoding/json"
	"fmt"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().StringP("env", "E", "dev", "env")
	updateCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
	updateCmd.Flags().StringP("name", "A", "", "name (required)")
	updateCmd.Flags().StringP("sql", "S", "postgres://postgres:password@192.168.1.7:5432", "sql url")
	updateCmd.MarkFlagRequired("name")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update node",
	Long:  "update node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		name, _ := cmd.Flags().GetString("name")
		sql, _ := cmd.Flags().GetString("sql")

		etcds := []string{host}
		prefix := fmt.Sprintf("%s/node", env)
		cx := core.EtcdAtomic{Endpoints: etcds}
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			v, err := ctx.Get(name)
			if err != nil {
				return err
			}
			cnf := conf.Config{}
			err = json.Unmarshal([]byte(v), &cnf)
			if err != nil {
				return err
			}
			cnf.DatabaseURL = sql
			return ctx.Put(name, string(util.ToJson(cnf)))
		})
		if err != nil {
			fmt.Printf("view command failed %s\n", err.Error())
		}
	},
}
