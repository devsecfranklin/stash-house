# README.md

[![DOI](https://zenodo.org/badge/407849291.svg)](https://zenodo.org/badge/latestdoi/407849291)

## Golang

```sh
go mod init github.com/devsecfranklin/website
go mod tidy
go mod verify
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
