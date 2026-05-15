package util

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Host       string
	User       string
	Password   string
	PrivateKey string
	PublicKey  string
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

	signer, err := ssh.ParsePrivateKey([]byte(c.PrivateKey))
	if err != nil {
		return err
	}
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(c.PublicKey)))
	if err != nil {
		fmt.Printf("%s\n", c.PublicKey)
		return err
	}
	conf := ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(pk),
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

func (c *SshClient) Upload(f os.File, p string, m string) error {
	cp, err := scp.NewClientBySSH(c.conn)
	if err != nil {
		return err
	}
	return cp.CopyFromFile(context.Background(), f, p, m)
}

func (c *SshClient) Download(f *os.File, p string, m string) error {
	cp, err := scp.NewClientBySSH(c.conn)
	if err != nil {
		return err
	}
	return cp.CopyFromRemote(context.Background(), f, p)
}
