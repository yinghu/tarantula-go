package node

import (
	"encoding/json"
	"fmt"

	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	resetCmd.Flags().StringP("env", "E", "dev", "env")
	resetCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
	resetCmd.Flags().StringP("app", "A", "", "app (required)")
	resetCmd.MarkFlagRequired("app")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset node",
	Long:  "reset node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		app, _ := cmd.Flags().GetString("app")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		prefix := fmt.Sprintf("%s/node", env)
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			clist := make([]conf.Config, 0)
			err := ctx.List(app, func(k, v string) bool {
				fmt.Printf("%s : %s\n", k, v)
				c := conf.Config{}
				err := json.Unmarshal([]byte(v), &c)
				if err != nil {
					return true
				}
				clist = append(clist, c)
				return true
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("reset command failed %s\n", err.Error())
		}
	},
}
