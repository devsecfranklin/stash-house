# docs

- Add dopcs here in markdown format.
- Convert to TeX, HTML, etc. as needed.

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
