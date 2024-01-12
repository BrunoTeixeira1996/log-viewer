package target

import (
	"net"
	"strings"
)

type Target struct {
	Name   string
	Host   string
	Status string
}

func (t *Target) IsListening(ip string) {
	_, err := net.Dial("tcp", ip+":9292")
	if err != nil {
		if strings.Contains(err.Error(), "connect") {
			t.Status = "unavailable"
			return
		}
	}
	t.Status = "available"
}
