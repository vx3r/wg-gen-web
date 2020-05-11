# Wg Gen Web

<h1 align="center"><img height="420" src="./wg-gen-web_cover.png" alt="Simple Web based configuration generator for WireGuard"></h1>

Simple Web based configuration generator for [WireGuard](https://wireguard.com).

[![pipeline status](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/badges/master/pipeline.svg)](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/commits/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/vx3r/wg-gen-web)](https://goreportcard.com/report/github.com/vx3r/wg-gen-web)
![Gitlab pipeline status (self-hosted)](https://img.shields.io/gitlab/pipeline/vx3r/wg-gen-web?gitlab_url=https%3A%2F%2Fgitlab.127-0-0-1.fr%2F)
[![License: WTFPL](https://img.shields.io/badge/License-WTFPL-brightgreen.svg)](http://www.wtfpl.net/about/)
![Discord](https://img.shields.io/discord/681699554189377567)
![Build multi-arch Docker Images via buildx](https://github.com/vx3r/wg-gen-web/workflows/Build%20multi-arch%20Docker%20Images%20via%20buildx/badge.svg)
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
 * Sent email to client with QR-code and client config
 * Enable / Disable client
 * Generation of `wg0.conf` after any modification
 * IPv6 ready
 * User authentication (Oauth2 OIDC)
 * Dockerized
 * Pretty cool look

![Screenshot](wg-gen-web_screenshot.png)

## Running

### Docker

The easiest way to run Wg Gen Web is using the container image
```
docker run --rm -it -v /tmp/wireguard:/data -p 8080:8080 -e "WG_CONF_DIR=/data" vx3r/wg-gen-web:latest
```
Docker compose snippet, used for demo server
```
version: '3.6'
  wg-gen-web-demo:
    image: vx3r/wg-gen-web:latest
    container_name: wg-gen-web-demo
    restart: unless-stopped
    expose:
      - "8080/tcp"
    environment:
      - WG_CONF_DIR=/data
      - WG_INTERFACE_NAME=wg0.conf
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USERNAME=no-reply@gmail.com
      - SMTP_PASSWORD=******************
      - SMTP_FROM=Wg Gen Web <no-reply@gmail.com>
      - OAUTH2_PROVIDER_NAME=github
      - OAUTH2_PROVIDER=https://github.com
      - OAUTH2_CLIENT_ID=******************
      - OAUTH2_CLIENT_SECRET=******************
      - OAUTH2_REDIRECT_URL=https://wg-gen-web-demo.127-0-0-1.fr
    volumes:
      - /etc/wireguard:/data
```
Please note that mapping ```/etc/wireguard``` to ```/data``` inside the docker, will erase your host's current configuration.
If needed, please make sure to backup your files from ```/etc/wireguard```.

A workaround would be to change the ```WG_INTERFACE_NAME``` to something different, as it will create a new interface (```wg-auto.conf``` for example), note that if you do so, you will have to adapt your daemon accordingly.

### Directly without docker

Fill free to download latest artifacts from my GitLab server:
* [Backend](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/-/jobs/artifacts/master/download?job=build-back)
* [Frontend](https://gitlab.127-0-0-1.fr/vx3r/wg-gen-web/-/jobs/artifacts/master/download?job=build-front)

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

## Authentication

Wg Gen Web can use Oauth2 OpenID Connect provider to authenticate users.
Currently there are 4 implementations:
- `fake` not a real implementation, use this if you don't want to authenticate your clients.

Add the environment variable:

```
OAUTH2_PROVIDER_NAME=fake
```

- `github` in order to use GitHub as Oauth2 provider.

Add the environment variable:

```
OAUTH2_PROVIDER_NAME=github
OAUTH2_PROVIDER=https://github.com
OAUTH2_CLIENT_ID=********************
OAUTH2_CLIENT_SECRET=********************
OAUTH2_REDIRECT_URL=https://wg-gen-web-demo.127-0-0-1.fr
```

- `google` in order to use Google as Oauth2 provider. Not yet implemented
```
help wanted
```

- `oauth2oidc` in order to use RFC compliant Oauth2 OpenId Connect provider.

Add the environment variable:

```
OAUTH2_PROVIDER_NAME=oauth2oidc
OAUTH2_PROVIDER=https://gitlab.com
OAUTH2_CLIENT_ID=********************
OAUTH2_CLIENT_SECRET=********************
OAUTH2_REDIRECT_URL=https://wg-gen-web-demo.127-0-0-1.fr
```

Please fell free to test and report any bugs.
Wg Gen Web will only access your profile to get email address and your name, no other unnecessary scopes will be requested.

## Need Help

 * Join us on [Discord](https://discord.gg/fjx7gGJ)
 * Create an issue

## Development

### Backend

From the top level directory run

```
$ go run main.go
```

### Frontend

Inside another terminal session navigate into the `ui` folder

```
$ cd ui
```
Install required dependencies
```
$ npm install
```
Set the base url for the api
```
$ export VUE_APP_API_BASE_URL=http://localhost:8080/api/v1.0
```
Start the development server. It will rebuild and reload the site once you make a change to the source code.
```
$ npm run serve
```

Now you can access the site from a webbrowser with the url `http://localhost:8081`.

## Application stack

 * [Gin, HTTP web framework written in Go](https://github.com/gin-gonic/gin)
 * [go-template, data-driven templates for generating textual output](https://golang.org/pkg/text/template/)
 * [Vue.js, progressive javaScript framework](https://github.com/vuejs/vue)
 * [Vuetify, material design component framework](https://github.com/vuetifyjs/vuetify)

## License

 * Do What the Fuck You Want to Public License. [LICENSE-WTFPL](LICENSE-WTFPL) or http://www.wtfpl.net
