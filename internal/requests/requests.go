package requests

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetJournalctlForTarget(targetHost string, targetName string) (interface{}, error) {
	// Perform request to extract info from journalctl
	req, err := http.NewRequest(http.MethodGet, "http://"+targetHost+":9292/log", nil)
	if err != nil {
		fmt.Printf("[ERROR] while creating new request: %s\n", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] while making http request: %s\n", err)
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("[ERROR] while reading response body: %s\n", err)
		return nil, err
	}

	data := struct {
		TargetHost string
		TargetName string
		Journalctl string
	}{
		TargetHost: targetHost,
		TargetName: targetName,
		Journalctl: strings.Replace(string(resBody), "\n", "<br>", -1),
	}

	return data, err
}
