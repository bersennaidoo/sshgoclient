package devops

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	expect "github.com/google/goexpect"
	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
)

type DevOps struct {
	ClientSSH *ssh.Client
}

func New(client *ssh.Client) *DevOps {
	return &DevOps{
		ClientSSH: client,
	}
}

func (dops *DevOps) SudoAptInstall(cmd string) (err error) {
	//timeout := 10 * time.Minute
	r, w := io.Pipe()
	debug := strings.Builder{}
	debugDone := make(chan struct{})
	go func() {
		io.Copy(&debug, r)
		close(debugDone)
	}()

	defer func() {
		// Wait for our io.Copy() to be done.
		<-debugDone

		// Only log this if we had an error.
		if err != nil {
			log.Printf("expect debug:\n%s", debug.String())
		}
	}()

	e, _, err := expect.SpawnSSH(dops.ClientSSH, 5*time.Second, expect.Tee(w))
	if err != nil {
		return err
	}
	defer e.Close()

	var promptRE = regexp.MustCompile(`# `)

	_, _, err = e.Expect(promptRE, 10*time.Second)
	if err != nil {
		return fmt.Errorf("did not get shell prompt")
	}

	if err := e.Send(fmt.Sprintf("sudo apt-get install %s -y\n", cmd)); err != nil {
		return fmt.Errorf("error on send command: %s", err)
	}

	e.ExpectBatch([]expect.Batcher{
		&expect.BCas{[]expect.Caser{
			&expect.Case{R: regexp.MustCompile(`[sudo] password for test: `), S: "test"},
		}},
	}, 10*time.Second)

	fmt.Println(term.Greenf("All done"))

	return nil
}

func (dops *DevOps) CombinedOutput(ctx context.Context, cmd string) (string, error) {
	sess, err := dops.ClientSSH.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	if v, ok := ctx.Deadline(); ok {
		t := time.NewTimer(v.Sub(time.Now()))
		defer t.Stop()

		go func() {
			x := <-t.C
			if !x.IsZero() {
				sess.Signal(ssh.SIGKILL)
			}
		}()
	}

	b, err := sess.Output(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
