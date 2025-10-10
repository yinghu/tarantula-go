package admin

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

func init() {
	sshCmd.Flags().StringP("host", "H", "192.168.1.7:22", "ssh host")
	sshCmd.Flags().StringP("user", "U", "", "user (required)")
	sshCmd.MarkFlagRequired("user")
	sshCmd.Flags().StringP("password", "P", "", "password (required)")
	sshCmd.MarkFlagRequired("password")
}

var (
	sshCmd = &cobra.Command{
		Use:   "ssh",
		Short: "ssh node",
		Long:  "ssh node",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			user, _ := cmd.Flags().GetString("user")
			password, _ := cmd.Flags().GetString("password")
			config := &ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
					ssh.Password(password),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			client, err := ssh.Dial("tcp", host, config)
			if err != nil {
				fmt.Printf("Failed to dial : %s\n", err)
			}
			session, err := client.NewSession()
			if err != nil {
				fmt.Printf("Failed to create session: %s\n", err)
			}
			defer session.Close()
			fmt.Printf("success!")
			modes := ssh.TerminalModes{
				ssh.ECHO:          0,     // disable echoing
				ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
				ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}
			if err := session.RequestPty("linux", 80, 40, modes); err != nil {
				fmt.Printf("request for pseudo terminal failed :%s\n", err)
			}

			// set input and output
			session.Stdout = os.Stdout
			session.Stdin = os.Stdin
			session.Stderr = os.Stderr

			if err := session.Shell(); err != nil {
				fmt.Printf("failed to start shell: %s\n", err)
			}
			fmt.Printf("success!")
			err = session.Wait()
			if err != nil {
				fmt.Printf("Failed to run: %s\n" ,err.Error())
			}
		},
	}
)
