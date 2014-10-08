package ovs

import (
	"errors"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type Flow struct {
	Table       int
	Duration    float64
	NumPackets  int
	NumBytes    int
	IdleAge     int
	Description map[string]string
}

type SwitchInfo struct {
	NumTables    int
	NumBuffers   int
	Capabilities []string
	Actions      []string
	Flows        []Flow
	LocalAddr    net.HardwareAddr
}

func FlowInfo(name string) (SwitchInfo, error) {
	var re *regexp.Regexp
	var matches []string
	var err error
	info := SwitchInfo{}

	// ovs-ofctl show
	stdout, stderr, err := run("ovs-ofctl", "show", name)
	if err != nil {
		errmsg := stderr.String()
		return info, errors.New(errmsg)
	}

	output := stdout.String()

	re = regexp.MustCompile("n_tables:(\\d+), n_buffers:(\\d+)")
	matches = re.FindStringSubmatch(output)
	if info.NumTables, err = strconv.Atoi(matches[1]); err != nil {
		return info, err
	}
	if info.NumBuffers, err = strconv.Atoi(matches[2]); err != nil {
		return info, err
	}

	re = regexp.MustCompile("capabilities: (.*)")
	matches = re.FindStringSubmatch(output)
	info.Capabilities = strings.Split(matches[1], " ")

	re = regexp.MustCompile("actions: (.*)")
	matches = re.FindStringSubmatch(output)
	info.Actions = strings.Split(matches[1], " ")

	re = regexp.MustCompile("LOCAL\\(.*?\\): addr:([0-9a-f\\:]+)")
	matches = re.FindStringSubmatch(output)
	if info.LocalAddr, err = net.ParseMAC(matches[1]); err != nil {
		return info, err
	}

	// ovs-ofctl dump-flows
	stdout, stderr, err = run("ovs-ofctl", "dump-flows", name)
	if err != nil {
		errmsg := stderr.String()
		return info, errors.New(errmsg)
	}

	for {
		line, err := stdout.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			return info, err
		}

		flow := Flow{}

		re = regexp.MustCompile("cookie=.*?, duration=([\\d\\.]+)s, table=(\\d+), n_packets=(\\d+), n_bytes=(\\d+), idle_age=(\\d+), (.*)")

		if !re.MatchString(line) {
			continue
		}

		matches = re.FindStringSubmatch(line)
		if flow.Duration, err = strconv.ParseFloat(matches[1], 64); err != nil {
			return info, err
		}
		if flow.Table, err = strconv.Atoi(matches[2]); err != nil {
			return info, err
		}
		if flow.NumPackets, err = strconv.Atoi(matches[3]); err != nil {
			return info, err
		}
		if flow.NumBytes, err = strconv.Atoi(matches[4]); err != nil {
			return info, err
		}
		if flow.IdleAge, err = strconv.Atoi(matches[5]); err != nil {
			return info, err
		}

		re = regexp.MustCompile("([^\\s]+)\\=([^\\s]+)")
		for _, desc := range re.FindAllStringSubmatch(matches[6], 0) {
			flow.Description[desc[1]] = desc[2]
		}

		info.Flows = append(info.Flows, flow)
	}

	return info, nil
}

func AddFlow(switch_name string, flow string) error {
	_, stderr, err := run("ovs-ofctl", "add-flow", switch_name, flow)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}

func DeleteFlow(switch_name string, flow string) error {
	_, stderr, err := run("ovs-ofctl", "del-flows", switch_name, flow)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}

func DeleteAllFlows(switch_name string) error {
	_, stderr, err := run("ovs-ofctl", "del-flows", switch_name)
	if err != nil {
		errmsg := stderr.String()
		return errors.New(errmsg)
	}

	return nil
}
