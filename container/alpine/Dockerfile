FROM alpine:3.20

RUN apk add --no-cache pass

RUN addgroup -g 1337 house && adduser -u 9001 -G house -s /bin/ash -h /home/stash -D stash
COPY container/alpine/stash/.ashrc /home/stash
COPY container/alpine/stash/.profile /home/stash
# RUN chown -R 9001:1337 /home/stash

CMD ["/bin/ash"]
