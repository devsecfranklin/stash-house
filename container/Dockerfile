# syntax=docker/dockerfile:1

FROM nginx:latest

LABEL org.opencontainers.image.authors="Franklin <franklin@dead10c5.org>"
LABEL org.opencontainers.image.source="https://github.com/devsecfranklin/website-bitsmasher.net"
LABEL org.opencontainers.image.description="bitsmasher.net website"
LABEL org.opencontainers.image.licenses=MIT

COPY website/assets /usr/share/nginx/html/assets
#COPY webstie/assets/js/env.arm.js /usr/share/nginx/html/assets/js/env.js
COPY package.json /usr/share/nginx/html
# COPY . /usr/share/nginx/html
COPY website/ /usr/share/nginx/html/
#COPY website/404.html /usr/share/nginx/html/404.html


CMD ["nginx", "-g", "daemon off;"]