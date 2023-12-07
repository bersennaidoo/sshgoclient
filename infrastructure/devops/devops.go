package devops

import "golang.org/x/crypto/ssh"

type DevOps struct {
	sclient *ssh.Client
}

func New(client *ssh.Client) *DevOps {
	return &DevOps{
		sclient: client,
	}
}
