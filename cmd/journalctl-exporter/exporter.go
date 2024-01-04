package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

/*
   Listen on port 9292
   Runs journalctl <service> like ./exporter -service gbackup.service
   Returns journalctl to http://localhost:9292/log
*/

type HandlerData struct {
	service string
}

// Handles "/log"
func (h *HandlerData) logHandle(w http.ResponseWriter, r *http.Request) {
	cmd, err := exec.Command("sudo", "journalctl", "-b", "-u", h.service, "-o", "short-precise").Output()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, string(cmd))
	}
}

func run() error {
	var serviceFlag = flag.String("service", "", "service to expose")
	flag.Parse()

	if *serviceFlag == "" {
		return fmt.Errorf("[ERROR] service flag is empty")
	}

	h := &HandlerData{service: *serviceFlag}

	mux := http.NewServeMux()
	mux.HandleFunc("/log", h.logHandle)
	http.ListenAndServe(":9292", mux)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
