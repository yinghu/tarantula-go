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
	updateCmd.Flags().String("etcd", "192.168.1.7:2379", "etcd host")
	updateCmd.Flags().StringP("app", "A", "", "app (required)")
	updateCmd.Flags().StringP("sql", "S", "postgres://postgres:password@192.168.1.7:5432", "sql url")
	updateCmd.Flags().StringP("http", "H", "192.168.1.11", "http host")
	updateCmd.Flags().StringP("tcp", "T", "tcp://192.168.1.11:5050", "tcp url")
	updateCmd.MarkFlagRequired("app")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update node",
	Long:  "update node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("etcd")
		app, _ := cmd.Flags().GetString("app")
		sql, _ := cmd.Flags().GetString("sql")
		http, _ := cmd.Flags().GetString("http")
		tcp, _ := cmd.Flags().GetString("tcp")
		etcds := []string{host}
		prefix := fmt.Sprintf("%s/node", env)
		cx := core.EtcdAtomic{Endpoints: etcds}
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			v, err := ctx.Get(app)
			if err != nil {
				return err
			}
			cnf := conf.Config{}
			err = json.Unmarshal([]byte(v), &cnf)
			if err != nil {
				return err
			}
			cnf.SqlEndpoint = sql
			cnf.HttpEndpoint = http
			cnf.TcpEndpoint = tcp
			return ctx.Put(app, string(util.ToJson(cnf)))
		})
		if err != nil {
			fmt.Printf("view command failed %s\n", err.Error())
		}
	},
}
