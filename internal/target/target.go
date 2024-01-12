package target

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"
)

type Target struct {
	Name   string
	Host   string
	Status string
}

func nmapForGokrazy(ip string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(ctx, nmap.WithTargets(ip), nmap.WithPorts("9292"))

	if err != nil {
		return false, fmt.Errorf("[ERROR] unable to create nmap scanner: %v", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		return false, fmt.Errorf("[ERROR] unable to run nmap scan: %v", err)
	}

	// Check if Host is up
	if result.Stats.Hosts.Up == 0 {
		return false, nil
	}

	// Check if port 9292 is open
	for _, port := range result.Hosts[0].Ports {
		if port.ID == 9292 {
			return true, nil
		}
	}

	return false, nil
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
