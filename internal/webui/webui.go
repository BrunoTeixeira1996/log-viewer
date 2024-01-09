package webui

import (
	"embed"
	"fmt"
	"log-viewer/internal/requests"
	"log-viewer/internal/target"
	"net/http"
	"strings"
	"text/template"
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

//go:embed assets/*
var assetsDir embed.FS

func Init(targets []target.Target) error {
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

	err = http.ListenAndServe(":8080", mux)
	if err != nil && err != http.ErrServerClosed {
		panic("Error trying to start http server: " + err.Error())
	}

	return nil
}
