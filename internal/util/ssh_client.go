package util

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Host     string
	User     string
	Password string
	conn     *ssh.Client
}

func (c *SshClient) Connect() error {
	conf := ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ci, err := ssh.Dial("tcp", c.Host, &conf)
	if err != nil {
		return err
	}
	c.conn = ci
	return nil
}

func (c *SshClient) Close() error {
	return c.conn.Close()
}

func (c *SshClient) Run(cmd string) error {
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	var buff bytes.Buffer
	session.Stdout = &buff
	session.Run(cmd)
	fmt.Printf("output %s\n", buff.String())
	return nil
}
