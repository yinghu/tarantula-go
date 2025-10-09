package node

import (
	"fmt"

	"gameclustering.com/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().StringP("env", "E", "dev", "env")
	cleanCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean nodes",
	Long:  "clean nodes",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		prefix := fmt.Sprintf("%s/node", env)
		err := cx.Execute(prefix, func(ctx core.Ctx) error {
			err := ctx.List("", func(k, v string) bool {
				fmt.Printf("%s : %s\n", k, v)
				ctx.Del(k, false)
				return true
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("clean command failed %s\n", err.Error())
		}
	},
}
