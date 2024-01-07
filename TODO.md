# TODO

- [X] Add `exporter.go` to journalctl-exporter that exposes a service
- [X] Expose journalctl service log in port `9292` in `exporter.go`
- [X] Add and parse toml file to `viewer.go`
- [X] Create webui that displays all targets, similar to gokrazy web interface
- [X] Check if target is alive when starting program
- [X] Create target webui similar to gokrazy when viewing logs of a specific appliance
- [ ] Perform request to `http://target_ip:9292/log` and display in web interface
