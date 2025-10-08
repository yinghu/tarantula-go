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
	addCmd.Flags().String("etcd", "192.168.1.7:2379", "etcd host")
	addCmd.Flags().StringP("app", "A", "", "app (required)")
	addCmd.Flags().StringP("sql", "S", "postgres://postgres:password@192.168.1.7:5432", "sql url")
	addCmd.Flags().StringP("http", "H", "192.168.1.11", "http host")
	addCmd.Flags().StringP("tcp", "T", "192.168.1.11", "tcp host")
	addCmd.MarkFlagRequired("app")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add node",
	Long:  "add node",
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
			nidKey := fmt.Sprintf("app.%s.nid", app)
			nid := 0
			nv, err := ctx.Get(nidKey)
			if err == nil {
				ic, err := strconv.ParseInt(nv, 10, 32)
				if err != nil {
					return err
				}
				nid = int(ic)
			}
			nv, err = ctx.Get("id")
			id := 0
			if err == nil {
				ic, err := strconv.ParseInt(nv, 10, 32)
				if err != nil {
					return err
				}
				id = int(ic)
			}
			if id > 1023 {
				return fmt.Errorf("snowflake id sequence must be 0 - 1023 %d", id)
			}
			cnf := conf.Config{Sequence: id}
			cnf.SqlEndpoint = sql
			cnf.HttpEndpoint = http
			cnf.TcpEndpoint = tcp
			name := fmt.Sprintf("%s.%d", app, nid)
			err = ctx.Put(name, string(util.ToJson(cnf)))
			if err != nil {
				return err
			}
			nid++
			id++

			err = ctx.Put(nidKey, strconv.Itoa(nid))
			if err != nil {
				return err
			}
			err = ctx.Put("id", strconv.Itoa(id))
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("failed to add node %s\n", err.Error())
			return
		}
		fmt.Printf("add node config %s on env %s\n", app, env)
	},
}
