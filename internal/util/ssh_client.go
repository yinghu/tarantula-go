package util

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Host       string
	User       string
	Password   string
	PrivateKey string
	conn       *ssh.Client
}

func (c *SshClient) WithPassword() error {
	conf := ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ci, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, 22), &conf)
	if err != nil {
		return err
	}
	c.conn = ci
	return nil
}

func (c *SshClient) WithKey() error {
	key, err := os.ReadFile(c.PrivateKey)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	conf := ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ci, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, 22), &conf)
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
	err = session.Run(cmd)
	if err != nil {
		return err
	}
	fmt.Printf("output %s\n", buff.String())
	return nil
}
