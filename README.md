# README.md

[![DOI](https://zenodo.org/badge/407849291.svg)](https://zenodo.org/badge/latestdoi/407849291) [![Cloudbot](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/cloudbot-call.yml/badge.svg)](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/cloudbot-call.yml) [![npm audit](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/npm-audit.yml/badge.svg)](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/npm-audit.yml)

## setup

Game and minecraft server

```sh
sudo apt -y install certbot python3-certbot-nginx
sudo certbot --nginx -d games.bitsmasher.net
certbot certificates # displays information about the certificates that Certbot has obtained
systemctl status certbot.timer # showing whether it’s active and when it’s scheduled to run next.
certbot renew --dry-run # If dry run is successful, the auto-renewal has been set up correctly.
```

## test web site locally

```bash
newgrp docker && bash
docker pull nginx
docker build -t nginx-bitsmasher .
docker run --name docker-nginx-bitsmasher -p 8080:80 -d nginx-bitsmasher
```

Now navigate to [http://0.0.0.0:8080/](http://0.0.0.0:8080/)

## test twitter card

[twitter card validator](https://cards-dev.twitter.com/validator)

## Linter

```sh
nix-shell shell.nix
linthtml **/*.html
for FILE in **/*.html; html-beautify {$FILE} > {$FILE}.tmp && mv {$FILE}.tmp {$FILE}; end
```

## Shell

```sh
apt install -y python3-pip screen
```
