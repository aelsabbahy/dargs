package run

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type CommandRunner struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func (c *CommandRunner) Key() string {
	return fmt.Sprint(strings.Join([]string{c.Name, c.Command}, "-"))
}

type EmptyOutError struct {
	cmd, serr string
}

func (e *EmptyOutError) Error() string {
	em := "STDOUT was empty for command"
	em += fmt.Sprint("\nCommand:\n", e.cmd)
	em += fmt.Sprintln("\n\nSTDERR:", e.serr)
	return em
}

func (c *CommandRunner) Run(arg string, match map[string]string) (string, error) {
	log := log.WithField("prefix", c.Name)
	log.Debug("Running filter")
	if c.Command == "" {
		e := "No command defined"
		log.Errorf(e)
		return "", fmt.Errorf(e)
	}
	var matchEnv []string
	for k, v := range match {
		matchEnv = append(matchEnv, fmt.Sprintf("RE_%s=%s", k, v))

	}
	sort.Strings(matchEnv)
	log.Debug("Match Vars:\n", strings.Join(matchEnv, "\n"))
	environ := append(os.Environ(), matchEnv...)
	run := c.Command
	log.Debug("Command:\n", run)
	sout, serr, err := runCmd(run, environ, "-e")
	if err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	if sout == "" {
		err = &EmptyOutError{
			cmd:  run,
			serr: serr,
		}
		return "", err
	}
	log.Debug("STDOUT:", sout)
	log.Debug("STDERR:", serr)
	return sout, nil
}
