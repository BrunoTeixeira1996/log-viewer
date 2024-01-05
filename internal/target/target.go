package target

import (
	"fmt"
	"os/exec"
	"strings"
)

type Target struct {
	Name   string
	Host   string
	Status string
}

func (t *Target) IsAlive(ip string) error {
	out, err := exec.Command("ping", ip, "-c 2").Output()
	if err != nil {
		t.Status = "unavailable"
		return fmt.Errorf("Could not ping that IP: %v", err)
	}

	if strings.Contains(string(out), "Destination Host Unreachable") {
		t.Status = "unavailable"
		return fmt.Errorf("Destination Host Unreachable")
	} else {
		t.Status = "available"
		return nil
	}
}
