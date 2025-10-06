package node

import (
	"fmt"
	"strconv"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().StringP("env", "E", "dev", "env")
	addCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
	addCmd.Flags().StringP("app", "A", "", "app (required)")
	addCmd.Flags().StringP("sql", "S", "postgres://postgres:password@192.168.1.7:5432", "sql url")
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
		sql, _ := cmd.Flags().GetString("sql")

		etcds := []string{host}
		prefix := fmt.Sprintf("%s/node", env)
		cx := core.EtcdAtomic{Endpoints: etcds}
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			sid, err := ctx.Get("id")
			id := 0
			if err == nil {
				i64, err := strconv.ParseInt(sid, 10, 32)
				if err != nil {
					return err
				}
				id = int(i64)
			}
			cnf := conf.Config{Sequence: id}
			cnf.DatabaseURL = sql
			name := fmt.Sprintf("%s.%d", app, cnf.Sequence)
			err = ctx.Put(name, string(util.ToJson(cnf)))
			if err != nil {
				return err
			}
			id++
			return ctx.Put("id", strconv.Itoa(id))
		})
		if err != nil {
			fmt.Printf("failed to add node %s\n", err.Error())
			return
		}
		fmt.Printf("add node %s %s %s\n", prefix, env, host)
	},
}
