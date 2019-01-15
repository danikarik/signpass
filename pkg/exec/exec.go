package exec

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

// Command is a wrapper around std exec cmd.
type Command struct {
	base   *exec.Cmd
	Name   string
	Args   []string
	Stderr io.ReadCloser
	Stdout io.ReadCloser
}

// New creates instance of Command struct.
func New(name string, args ...string) (*Command, error) {
	cmd := exec.Command(name, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("could not create stderr pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not create stdout pipe: %v", err)
	}
	return &Command{
		base:   cmd,
		Name:   name,
		Args:   args,
		Stderr: stderr,
		Stdout: stdout,
	}, nil
}

// Start executes command's script.
func (c *Command) Start() error {
	err := c.base.Start()
	if err != nil {
		return err
	}
	return nil
}

// Wait process finish.
func (c *Command) Wait() error {
	err := c.base.Wait()
	if err != nil {
		return err
	}
	return nil
}

// StderrStr returns stderr buffer in string.
func (c *Command) StderrStr() string {
	buf, err := ioutil.ReadAll(c.Stderr)
	if err != nil {
		fmt.Printf("[WARNING] stderr to string: %s\n", err.Error())
		return ""
	}
	cout := strings.Trim(strings.Trim(string(buf), "\n"), " ")
	return cout
}

// StdoutStr returns stdout buffer in string.
func (c *Command) StdoutStr() string {
	buf, err := ioutil.ReadAll(c.Stdout)
	if err != nil {
		fmt.Printf("[WARNING] stdout to string: %s\n", err.Error())
		return ""
	}
	cout := strings.Trim(strings.Trim(string(buf), "\n"), " ")
	return cout
}
