package ovs_test

import (
	"fmt"
	"github.com/mistifyio/go-openvswitch"
	"testing"
)

func TestCommands(t *testing.T) {
	var err error

	// AddBridge
	if err = ovs.AddBridge("br0"); err != nil {
		fmt.Printf("Add bridge failed: %s", err.Error())
		t.FailNow()
	}

	// ListBridges
	bridges, err := ovs.ListBridges()

	if err != nil {
		fmt.Printf("List bridges failed: %s", err.Error())
		t.FailNow()
	}

	if len(bridges) != 1 {
		fmt.Printf("List bridges returned %d bridges (should be 1).", len(bridges))
		t.FailNow()
	}

	// AddPort
	if err = ovs.AddPort("br0", "eth1"); err != nil {
		fmt.Printf("Add port failed: %s", err.Error())
		t.FailNow()
	}

	// ListPorts
	ports, err := ovs.ListPorts("br0")

	if err != nil {
		fmt.Printf("List ports failed: %s", err.Error())
		t.FailNow()
	}

	if len(ports) != 1 {
		fmt.Printf("List ports returned no ports.")
		t.FailNow()
	}

	// DeletePort
	if err = ovs.DeletePort("br0", "eth1"); err != nil {
		fmt.Printf("Delete port failed: %s", err.Error())
		t.FailNow()
	}

	ports, err = ovs.ListPorts("br0")

	if len(ports) != 0 {
		fmt.Printf("%d ports exist on br0 after delete (should be 0)", len(ports))
		t.FailNow()
	}

	// DeleteBridge
	if err = ovs.DeleteBridge("br0"); err != nil {
		fmt.Printf("Delete bridge failed: %s", err.Error())
		t.FailNow()
	}

	bridges, err = ovs.ListBridges()

	if len(bridges) != 0 {
		fmt.Printf("%d bridges exist after delete (should be 0)", len(bridges))
		t.FailNow()
	}
}

func TestOpenFlow(t *testing.T) {
	var err error

	// AddBridge
	if err = ovs.AddBridge("br0"); err != nil {
		fmt.Printf("Add bridge failed: %s", err.Error())
		t.FailNow()
	}

	// FlowInfo
	info, err := ovs.FlowInfo("br0")

	if err != nil {
		fmt.Printf("OpenFlow show failed: %s", err.Error())
		t.FailNow()
	}

	if len(info.Flows) == 0 {
		fmt.Printf("OpenFlow show failed: no flows")
		t.FailNow()
	}

	// AddFlow
	if err = ovs.AddFlow("br0", "table=0, dl_src=01:00:00:00:00:00/01:00:00:00:00:00, actions=drop"); err != nil {
		fmt.Printf("AddFlow failed: %s", err.Error())
		t.FailNow()
	}

	// DeleteFlow
	if err = ovs.DeleteFlow("br0", "table=0, dl_src=01:00:00:00:00:00/01:00:00:00:00:00"); err != nil {
		fmt.Printf("DelFlow failed: %s", err.Error())
		t.FailNow()
	}

	// DeleteBridge
	if err = ovs.DeleteBridge("br0"); err != nil {
		fmt.Printf("Delete bridge failed: %s", err.Error())
		t.FailNow()
	}
}
