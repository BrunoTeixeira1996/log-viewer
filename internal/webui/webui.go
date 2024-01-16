package webui

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"text/template"

	"github.com/BrunoTeixeira1996/log-viewer/internal/requests"
	"github.com/BrunoTeixeira1996/log-viewer/internal/target"
)

type UI struct {
	tmpl    *template.Template
	targets []target.Target
}

func (ui *UI) indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := ui.tmpl.ExecuteTemplate(w, "index.html.tmpl", map[string]interface{}{
		"targets": ui.targets,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ui *UI) rawHandler(w http.ResponseWriter, r *http.Request) {
	// Verifies if target exist in order to not go to /raw/<target>
	var targetHost string
	targetName := strings.TrimPrefix(r.URL.Path, "/raw/")

	for _, t := range ui.targets {
		if t.Name == targetName {
			targetHost = t.Host
			break
		}
	}

	if targetHost == "" {
		fmt.Fprintf(w, "Target does not exist")
		return
	}

	// Performs request to extract info from journalctl in raw
	data, err := requests.GetJournalctlForTarget(targetHost, targetName)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}
	// This extracts Journalctl field from interface
	requestValue := reflect.ValueOf(data)
	fmt.Fprintf(w, "%s\n", strings.Replace(requestValue.FieldByName("Journalctl").String(), "<br>", "\n", -1))

	return
}

func (ui *UI) targetHandler(w http.ResponseWriter, r *http.Request) {
	// Verifies if target exist in order to not go to /target/asdasdasd
	var targetHost string
	targetName := strings.TrimPrefix(r.URL.Path, "/target/")

	for _, t := range ui.targets {
		if t.Name == targetName {
			targetHost = t.Host
			break
		}
	}

	if targetHost == "" {
		fmt.Fprintf(w, "Target does not exist")
		return
	}

	// Performs request to extract info from journalctl
	data, err := requests.GetJournalctlForTarget(targetHost, targetName)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	if err := ui.tmpl.ExecuteTemplate(w, "target.html.tmpl", map[string]interface{}{
		"data": data,
	}); err != nil {
		fmt.Printf("[ERROR] while executing template: %s\n", err)
		fmt.Fprintf(w, "[ERROR] while executing template: %s\n", err)
		return
	}
	return
}

// Handles GET to view targets status on demand
func (ui *UI) statusHandle(w http.ResponseWriter, r *http.Request) {
	var res string

	defer r.Body.Close()
	if r.Method != "GET" {
		http.Error(w, "NOT GET!", http.StatusBadRequest)
		return
	}

	log.Printf("GET request from %s to check if targets are available on demand\n", r.RemoteAddr)
	for i, t := range ui.targets {
		t.IsListening(t.Host)
		(*&ui.targets)[i] = t
		if t.Status == "available" {
			res += "Target " + t.Name + " is avaiable\n"
		}
	}

	fmt.Fprintf(w, res)
	return
}

//go:embed assets/*
var assetsDir embed.FS

func Init(targets []target.Target, listenPort string) error {
	tmpl, err := template.ParseFS(assetsDir, "assets/*.tmpl")
	if err != nil {
		return err
	}

	ui := &UI{
		tmpl:    tmpl,
		targets: targets,
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.FileServer(http.FS(assetsDir)))
	mux.HandleFunc("/", ui.indexHandler)
	mux.HandleFunc("/target/", ui.targetHandler)
	mux.HandleFunc("/status/", ui.statusHandle)
	mux.HandleFunc("/raw/", ui.rawHandler)

	log.Printf("Listening at :%s\n", listenPort)

	err = http.ListenAndServe(":"+listenPort, mux)
	if err != nil && err != http.ErrServerClosed {
		panic("Error trying to start http server: " + err.Error())
	}

	return nil
}
