# README.md

[![DOI](https://zenodo.org/badge/407849291.svg)](https://zenodo.org/badge/latestdoi/407849291) [![Go Report Card](https://goreportcard.com/badge/github.com/golang-standards/project-layout?style=flat-square)](https://goreportcard.com/report/github.com/golang-standards/project-layout) [![Documentation](https://github.com/devsecfranklin/website/actions/workflows/generate-docs.yml/badge.svg)](https://github.com/devsecfranklin/website/actions/workflows/generate-docs.yml)

## Golang

```sh
go mod init github.com/devsecfranklin/website
go mod tidy # create/update go.mod
go mod verify # create go.sum
go run cmd/games/*.go # run the application on port 8080
```

## Certbot

Game and minecraft server

```sh
sudo apt -y install certbot python3-certbot-nginx
sudo certbot --nginx -d games.bitsmasher.net
certbot certificates # displays information about the certificates that Certbot has obtained
systemctl status certbot.timer # showing whether it’s active and when it’s scheduled to run next.
certbot renew --dry-run # If dry run is successful, the auto-renewal has been set up correctly.
```
