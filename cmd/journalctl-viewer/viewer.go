package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/BrunoTeixeira1996/log-viewer/internal/config"
	"github.com/BrunoTeixeira1996/log-viewer/internal/target"
	"github.com/BrunoTeixeira1996/log-viewer/internal/webui"

	"github.com/jinzhu/copier"
)

/*
   Listen on port 9191
   Reads toml file that has exporters with name and URL (http://192.168.30.23:9090/log)
   Has webserver in / listing all entries of the toml file and every entry is an a tag
   When a tag is pressed, go to /<target> and show journalctl output like gokrazy does
   Runs in gokrazy so I can go to gokrazy:9191 and view all logs
*/

// Goroutine that check if target is still listening
func isTargetStillListening(targets *[]target.Target, timeToCheckListening int) {
	for {
		log.Printf("After %d minutes, Ill check if targets are still listening ...\n", timeToCheckListening)
		time.Sleep(time.Duration(timeToCheckListening) * time.Minute)

		for i, t := range *targets {
			t.IsListening(t.Host)
			log.Printf("Target %s (%s) is %s\n", t.Name, t.Host, t.Status)
			(*targets)[i] = t
		}
	}
}

func run() error {
	var (
		cfg                  config.Config
		listenPortFlag       = flag.String("listen-port", "9696", "listening port (default is 9696)")
		tomlPathFlag         = flag.String("toml-file", "", "path to toml file")
		isStillListeningFlag = flag.String("check-time", "", "time to see if target is still listening")
		timeToCheckListening int
		err                  error
	)
	flag.Parse()

	if _, err = strconv.Atoi(*listenPortFlag); err != nil {
		return fmt.Errorf("[ERROR] please provide a valid listen-port flag")
	}

	if *tomlPathFlag == "" {
		return fmt.Errorf("[ERROR] toml-file flag is empty")
	}

	if timeToCheckListening, err = strconv.Atoi(*isStillListeningFlag); err != nil {
		return fmt.Errorf("[ERROR] please provide an int type for check-time flag")
	}

	if cfg, err = config.ReadTomlFile(*tomlPathFlag); err != nil {
		log.Fatal(err)
	}

	// Copies config target to target.Target (object displayed in webui)
	targets := &[]target.Target{}
	copier.Copy(&targets, &cfg.Targets)

	var wg sync.WaitGroup
	wg.Add(len(*targets))
	for i, t := range *targets {
		go func(i int, t target.Target) {
			t.IsListening(t.Host)
			log.Printf("Target %s (%s) is %s\n", t.Name, t.Host, t.Status)
			(*targets)[i] = t
			wg.Done()
		}(i, t)

	}
	wg.Wait()

	go isTargetStillListening(targets, timeToCheckListening)

	if err := webui.Init(*targets, *listenPortFlag); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
	}
}
