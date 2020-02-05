# Wg Gen Web

<h1 align="center"><img height="420" src="./wg-gen-web_cover.png" alt="Simple Web based configuration generator for WireGuard"></h1>

Simple Web based configuration generator for [WireGuard](https://wireguard.com).

[![pipeline status](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/badges/master/pipeline.svg)](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/commits/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/vx3r/wg-gen-web)](https://goreportcard.com/report/github.com/vx3r/wg-gen-web)
![Gitlab pipeline status (self-hosted)](https://img.shields.io/gitlab/pipeline/vx3r/wg-gen-web?gitlab_url=https%3A%2F%2Fgitlab.127-0-0-1.fr%2F)
[![License: WTFPL](https://img.shields.io/badge/License-WTFPL-brightgreen.svg)](http://www.wtfpl.net/about/)
![GitHub last commit](https://img.shields.io/github/last-commit/vx3r/wg-gen-web)
![Docker Pulls](https://img.shields.io/docker/pulls/vx3r/wg-gen-web)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vx3r/wg-gen-web)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/vx3r/wg-gen-web)

## Why another one ?

All WireGuard UI implementations are trying to manage the service by applying configurations and creating network rules.
This implementation only generates configuration and its up to you to create network rules and apply configuration to WireGuard.
For example by monitoring generated directory with [inotifywait](https://github.com/inotify-tools/inotify-tools/wiki). 

The goal is to run Wg Gen Web in a container and WireGuard on host system.

## Features

 * Self-hosted and web based
 * Automatically select IP from the netowrk pool assigned to client
 * QR-Code for convenient mobile client configuration
 * Enable / Disable client
 * Generation of `wg0.conf` after any modification
 * Dockerized
 * Pretty cool look

![Screenshot](wg-gen-web_screenshot.png)

## Running

### Docker

The easiest way to run Wg Gen Web is using the container image
```
docker run --rm -it -v /tmp/wireguard:/data -p 8080:8080 -e "WG_CONF_DIR=/data" vx3r/wg-gen-web:latest
```
Docker compose snippet
```
version: '3.6'
services:
  wg-gen-web:
    image: vx3r/wg-gen-web:latest
    container_name: wg-gen-web
    restart: unless-stopped
    expose:
      - "8080/tcp"
    environment:
      - WG_CONF_DIR=/data
      - WG_INTERFACE_NAME=wg0.conf
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USERNAME=account@gmail.com
      - SMTP_PASSWORD="*************"
      - SMTP_FROM="Wg Gen Web <account@gmail.com>"
    volumes:
      - /etc/wireguard:/data
```
Please note that mapping ```/etc/wireguard``` to ```/data``` inside the docker, will erase your host's current configuration.
If needed, please make sure to backup your files from ```/etc/wireguard```.

A workaround would be to change the ```WG_INTERFACE_NAME``` to something different, as it will create a new interface (```wg-auto.conf``` for example), note that if you do so, you will have to adapt your daemon accordingly.

### Directly without docker

Fill free to download latest artefacts from my GitLab server:
* [Backend](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/-/jobs/artifacts/master/download?job=build-front)
* [Frontend](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/-/jobs/artifacts/master/download?job=build-back)

Put everything in one directory, create `.env` file with all configurations and run the backend.

## Automatically apply changes to WireGuard

### Using ```systemd```
Using `systemd.path` monitor for directory changes see [systemd doc](https://www.freedesktop.org/software/systemd/man/systemd.path.html)
```
# /etc/systemd/system/wg-gen-web.path
[Unit]
Description=Watch /etc/wireguard for changes

[Path]
PathModified=/etc/wireguard

[Install]
WantedBy=multi-user.target
```
This `.path` will activate unit file with the same name
```
# /etc/systemd/system/wg-gen-web.service
[Unit]
Description=Restart WireGuard
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl restart wg-quick@wg0.service

[Install]
WantedBy=multi-user.target
```
Which will restart WireGuard service 

### Using ```inotifywait```
For any other init system, create a daemon running this script
```
#!/bin/sh
while inotifywait -e modify -e create /etc/wireguard; do
  wg-quick down wg0
  wg-quick up wg0
done
```

## How to use with existing WireGuard configuration

After first run Wg Gen Web will create `server.json` in data directory with all server informations.

Feel free to modify this file in order to use your existing keys

## What is out of scope

 * Generation or application of any `iptables` or `nftables` rules
 * Application of configuration to WireGuard by Wg Gen Web itself

## TODO

 * Multi-user support behind [Authelia](https://github.com/authelia/authelia) (suggestions / thoughts are welcome)
 * ~~Send configs by email to client~~

## License

 * Do What the Fuck You Want to Public License. [LICENSE-WTFPL](LICENSE-WTFPL) or http://www.wtfpl.net