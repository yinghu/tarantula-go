package node

import (
	"fmt"

	"gameclustering.com/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	viewCmd.Flags().StringP("env", "E", "dev", "env")
	viewCmd.Flags().StringP("host", "H", "192.168.1.7:2379", "etcd host")
	viewCmd.Flags().StringP("app", "A", "", "app (required)")
	viewCmd.MarkFlagRequired("app")
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view node",
	Long:  "view node",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		host, _ := cmd.Flags().GetString("host")
		app, _ := cmd.Flags().GetString("app")
		etcds := []string{host}
		cx := core.EtcdAtomic{Endpoints: etcds}
		prefix := fmt.Sprintf("%s.%s", env, app)
		//cnf := conf.Config{Sequence: 1}
		cx.Execute(prefix, func(ctx core.Ctx) error {
			nds, err := ctx.List(app, func(k, v string) {
				fmt.Printf("%k : %v\n", k, v)
			})
			if err != nil {
				return err
			}
			fmt.Printf("LIST : %v\n", nds)
			return nil
		})

	},
}
