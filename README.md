# Description

log-viewer is a tool to centralize logs in one single location.

Most of my scripts run in a raspberry pi with gokrazy. However other scripts that I need to run are hard to implement in the same gokrazy instance and others just don't belong to that stack. For example, [gbackup](https://github.com/BrunoTeixeira1996/gbackup) is a backup utility that I run inside proxmox, so all logs are being thrown to journalctl. Another example is a telegram bot that I use to perform some tasks. This bot does not fit inside the gokrazy instance so I use that as a systemd service.

With that, this utility aims to have two parts, an `exporter` and a `viewer`. The `exporter` exposes the journalctl of a systemd service and the `viewer` consumes that information and displays it in a simple web app. This utility is aimed to execute on gokrazy, however it can be executed as a systemd service and even expose logs and consume it self.

# Installation

## Exporter

- Go to the server where you have your service running and download the `journalctl-exporter` binary from this repository
- Then execute the binary like the following `./journalctl-exporter -service <service>`

## Viewer

- Go to the centralized server (in my case I use gokrazy) and download the `journalctl-viewer` binary from this repository
- Create a `config.toml` file and add your targets (you can follow the example `config.toml` file inside this repository)
- Then execute the binary like the following `./journalctl-viewer -toml-file <config.toml> -check-time <timer to check if targets are still listening>`
- Note that you need `sudo` and `nmap` installed in the viewer server

# Screenshots

![image](https://github.com/BrunoTeixeira1996/log-viewer/assets/12052283/6cd8e30e-6a84-4c4a-9509-127e9fb68ace)


![image](https://github.com/BrunoTeixeira1996/log-viewer/assets/12052283/104948d7-61c8-407d-b09d-106957afbc3c)
