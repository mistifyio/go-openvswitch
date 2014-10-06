package ovs_test

import (
	"fmt"
	"github.com/mistifyio/go-openvswitch"
	"testing"
)

func TestAddBridge(t *testing.T) {
	if err := ovs.AddBridge("br0"); err != nil {
		fmt.Printf("Add bridge failed: %s", err.Error())
		t.FailNow()
	}
}
