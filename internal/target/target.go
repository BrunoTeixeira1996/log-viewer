package target

import (
	"os/exec"
	"strings"
)

type Target struct {
	Name   string
	Host   string
	Status string
}

func (t *Target) IsListening(ip string) {
	out, err := exec.Command("nmap", "-p", "9292", ip).Output()
	if err != nil {
		t.Status = "unavailable"
	}

	if strings.Contains(string(out), "open") {
		t.Status = "available"
	} else {
		t.Status = "unavailable"
	}
}
