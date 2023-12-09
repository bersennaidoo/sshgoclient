package main

import (
	"flag"
	"fmt"
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

	client := sshc.New(*private, "root", nargs[0])

	dops := devops.New(client)
	defer dops.ClientSSH.Close()

	_ = dops.SudoAptInstall("tree")
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()

	//out, _ := dops.CombinedOutput(ctx, "date")
	//fmt.Printf("The Date is: %s\n", out)
}
