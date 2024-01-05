package main

import (
	"flag"
	"fmt"
	"log"
	"log-viewer/internal/config"
	"log-viewer/internal/target"
	"log-viewer/internal/webui"

	"github.com/jinzhu/copier"
)

/*
   Listen on port 9191
   Reads toml file that has exporters with name and URL (http://192.168.30.23:9090/log)
   Has webserver in / listing all entries of the toml file and every entry is an a tag
   When a tag is pressed, go to /<target> and show journalctl output like gokrazy does
   Runs in gokrazy so I can go to gokrazy:9191 and view all logs
*/

func run() error {
	var (
		cfg          config.Config
		tomlPathFlag = flag.String("toml-file", "", "path to toml file")
		err          error
	)
	flag.Parse()

	if *tomlPathFlag == "" {
		return fmt.Errorf("[ERROR] toml-file flag is empty")
	}

	if cfg, err = config.ReadTomlFile(*tomlPathFlag); err != nil {
		log.Fatal(err)
	}

	// Copies config target to target.Target (object displayed in webui)
	targets := &[]target.Target{}
	copier.Copy(&targets, &cfg.Targets)
	for i, t := range *targets {
		if err := t.IsAlive(t.Host); err != nil {
			log.Printf("Error while checking if IP is alive %s %v", t.Host, err)
		}
		(*targets)[i] = t
	}

	if err := webui.Init(*targets); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
	}
}
