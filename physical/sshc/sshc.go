package sshc

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func New(private string, user string, host string) *ssh.Client {
	var auth ssh.AuthMethod

	if private == "" {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			fmt.Println("-private not set, cannot use password when STDIN is a pipe")
			os.Exit(1)
		}
	} else {
		auth, err = publicKey(private)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	u := user

	config := &ssh.ClientConfig{
		User:            u,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		fmt.Println("Error: could not dial host: ", err)
		os.Exit(1)
	}

	return client
}

func publicKey(privateKeyFile string) (ssh.AuthMethod, error) {
	k, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(k)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
