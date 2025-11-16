# Docker

## auth

```bash
ssh -T git@github.com
git config --global --add url."git@github.com:".insteadOf "https://github.com/"
```

## build and run

Fish shell:

```fish
docker build -t frank378:beamer-example --build-arg \
  BUILD_DATE=(date -u +'%Y-%m-%dT%H:%M:%SZ') .
```

```bash
docker build -t frank378:beamer-example \
--build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') .
```

```bash
#Successfully built cdf6b6fafe03
docker inspect cdf6b6fafe03
docker images
docker run --rm  -it --entrypoint /bin/bash cdf6b6fafe03
```

```bash
docker run -it --entrypoint /bin/bash -e "TERM=xterm-256color" cdf6b6fafe03
```
