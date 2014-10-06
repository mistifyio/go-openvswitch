// Package ovs provides wrappers around the Open vSwitch command-line tools
package ovs

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
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

func ListBridges() ([]string, error) {
	var bridges []string

	stdout, stderr, err := run("ovs-vsctl", "list-br")
	if err != nil {
		errmsg := stderr.String()
		return bridges, errors.New(errmsg)
	}

	bridges = strings.Split(stdout.String(), "\n")
	bridges = bridges[:len(bridges)-1] // Remove empty line from end
	return bridges, nil
}

func DeleteBridge(name string) error {
	_, stderr, err := run("ovs-vsctl", "del-br", name)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}

func AddPort(bridge string, port string) error {
	_, stderr, err := run("ovs-vsctl", "add-port", bridge, port)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}

func ListPorts(bridge string) ([]string, error) {
	var ports []string

	stdout, stderr, err := run("ovs-vsctl", "list-ports", bridge)
	if err != nil {
		errmsg := stderr.String()
		return ports, errors.New(errmsg)
	}

	ports = strings.Split(stdout.String(), "\n")
	ports = ports[:len(ports)-1] // Remove empty line from end
	return ports, nil
}

func DeletePort(bridge string, port string) error {
	_, stderr, err := run("ovs-vsctl", "del-port", bridge, port)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}
