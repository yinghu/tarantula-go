package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("user", "U", "", "user (required)")
	loginCmd.MarkFlagRequired("user")
	loginCmd.Flags().StringP("password", "P", "", "password (required)")
	loginCmd.MarkFlagRequired("password")
	loginCmd.Flags().StringP("host", "H", "http://localhost", "use localhost")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
	Long:  "login authentication",
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		host, _ := cmd.Flags().GetString("host")
		hc := util.HttpCaller{Host: host}
		login := bootstrap.Login{Name: user, Hash: password}
		err := hc.PostJson("presence/login", login, func(resp *http.Response) error {
			session := core.OnSession{}
			err := json.NewDecoder(resp.Body).Decode(&session)
			if err != nil {
				return err
			}
			if !session.Successful {
				return fmt.Errorf("error : %s", session.Message)
			}
			hc.Token = session.Token
			hc.Ticket = session.Ticket
			hc.SystemId = session.SystemId
			hc.Home = session.Home
			return nil
		})
		if err != nil {
			fmt.Printf("authencation failed %s\n", err.Error())
			return
		}
		fmt.Printf("Authenticated %d\n", hc.SystemId)
	},
}
