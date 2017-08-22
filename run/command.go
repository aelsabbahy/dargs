package run

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runCmd(run string, environ []string, shellArgs ...string) (string, string, error) {
	cmd := exec.Command("bash", shellArgs...)
	var outbuf, errbuf bytes.Buffer
	cmd.Env = environ
	cmd.Stdin = strings.NewReader(run)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	sout := strings.TrimSpace(outbuf.String())
	serr := strings.TrimSpace(errbuf.String())
	if err != nil {
		em := fmt.Sprintf("Command Failed with error: %s\n", err)
		em += fmt.Sprint("\nCommand:\n", run)
		em += fmt.Sprintln("\nSTDOUT:", sout)
		em += fmt.Sprintln("\nSTDERR:", serr)
		return sout, serr, fmt.Errorf(em)
	}
	return sout, serr, err
}
