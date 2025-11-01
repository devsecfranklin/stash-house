FROM golang:alpine

RUN apk add --no-cache pass

RUN addgroup -g 1337 house && adduser -u 9001 -G house -s /bin/ash -h /home/stash -D stash
COPY container/alpine/stash/.ashrc /home/stash
COPY container/alpine/stash/.profile /home/stash
# RUN chown -R 9001:1337 /home/stash

# sudo apt-get install gcc make libpcsclite-dev pkg-config
# yubikey support

# go install github.com/foks-proj/go-foks/client/foks@latest
# (cd $(dirname $(which foks)) && ln -s foks git-remote-foks)
# foks ctl start 

CMD ["/bin/ash"]
