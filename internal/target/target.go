package target

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

type Target struct {
	Name   string
	Host   string
	Status string
}

func nmapForGokrazy(ip string) (bool, error) {
	_, err := net.Dial("tcp", ip+":9292")
	// Is not up nor listening on port 9292
	if err != nil {
		if strings.Contains(err.Error(), "connect") {
			return false, nil
		}
	}
	return true, nil
}

func (t *Target) IsListening(ip string, isGokrazy bool) {
	switch isGokrazy {
	case false:
		out, err := exec.Command("nmap", "-p", "9292", ip).Output()
		if err != nil {
			t.Status = "unavailable"
		}

		if strings.Contains(string(out), "open") {
			t.Status = "available"
		} else {
			t.Status = "unavailable"
		}

	case true:
		isPortOpen, err := nmapForGokrazy(ip)
		if err != nil {
			log.Println("[ERROR] while doing nmapForGokrazy:", err)
			break
		}

		if isPortOpen {
			t.Status = "available"
		} else {
			t.Status = "unavailable"
		}

	}

}
