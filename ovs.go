// Package ovs provides wrappers around the Open vSwitch command-line tools
package ovs

import (
	"bytes"
	"errors"
	"os/exec"
)

func run(cmd string, arg ...string) (bytes.Buffer, bytes.Buffer, error) {
	var stdout, stderr bytes.Buffer

	command := exec.Command("ovs-vsctl", arg...)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	return stdout, stderr, err
}

func AddBridge(name string) error {
	_, stderr, err := run("ovs-vsctl", "add-br", name)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}
