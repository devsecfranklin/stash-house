# README.md

[![DOI](https://zenodo.org/badge/407849291.svg)](https://zenodo.org/badge/latestdoi/407849291) [![Cloudbot](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/cloudbot-call.yml/badge.svg)](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/cloudbot-call.yml) [![npm audit](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/npm-audit.yml/badge.svg)](https://github.com/devsecfranklin/website-bitsmasher.net/actions/workflows/npm-audit.yml)

## test web site locally

```bash
newgrp docker && bash
docker pull nginx
docker build -t nginx-bitsmasher .
docker run --name docker-nginx-bitsmasher -p 8080:80 -d nginx-bitsmasher
```

Now navigate to http://http://0.0.0.0:8080/

## test twitter card

[twitter card validator](https://cards-dev.twitter.com/validator)
