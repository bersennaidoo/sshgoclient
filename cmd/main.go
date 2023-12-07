package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bersennaidoo/sshgoclient/infrastructure/devops"
	"github.com/bersennaidoo/sshgoclient/physical/sshc"
)

var (
	private = flag.String("private", "", "The path to the SSH private key for this connection")
)

func main() {
	flag.Parse()

	nargs := flag.Args()

	if len(nargs) != 1 {
		fmt.Println("Error: command must be 1 arg, [host]")
		os.Exit(1)
	}
	_, _, err := net.SplitHostPort(nargs[0])
	if err != nil {
		nargs[0] = nargs[0] + ":22"
		_, _, err = net.SplitHostPort(os.Args[1])
		if err != nil {
			fmt.Println("Error: problem with host passed: ", err)
			os.Exit(1)
		}
	}

	client := sshc.New(*private, "test", nargs[0])

	dops := devops.New(client)
	log.Printf("%#q\n", dops)

}
