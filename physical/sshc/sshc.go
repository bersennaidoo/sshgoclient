package sshc

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func New(private string, user string, host string) *ssh.Client {
	var auth ssh.AuthMethod

	if private == "" {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			fmt.Println("-private not set, cannot use password when STDIN is a pipe")
			os.Exit(1)
		}
		auth1, err := passwordFromTerm()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		auth = auth1
	} else {
		auth2, err := publicKey(private)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		auth = auth2
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

func passwordFromTerm() (ssh.AuthMethod, error) {
	fmt.Printf("SSH Passsword: ")
	p, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	fmt.Println("") // Show the return
	if len(bytes.TrimSpace(p)) == 0 {
		return nil, fmt.Errorf("password was an empty string")
	}
	return ssh.Password(string(p)), nil
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
