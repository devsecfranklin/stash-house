# syntax=docker/dockerfile:1

FROM python:3.13-bullseye

LABEL org.opencontainers.image.authors="Franklin <franklin@bitsmasher.net>"
LABEL org.opencontainers.image.source="https://github.com/devsecfranklin/minecraft"
LABEL org.opencontainers.image.description="Minecraft server files and tools."
LABEL org.opencontainers.image.licenses=MIT

RUN mkdir -p /home/franklin/workspace
WORKDIR /home/franklin/workspace
ENV MY_DIR /home/franklin/workspace
COPY . ${MY_DIR}

# COPY docker/python3/files/bashrc /workspaces/franklin-resume/.bashrc

ENV DEBIAN_FRONTEND noninteractive

RUN \
  pip install --upgrade pip && \
  pip install Cython && \
  pip install -r${MY_DIR}/src/requirements.txt && \
  apt update && \
  apt install -y make git

CMD ["/usr/local/bin/python3", "app/app.py"]
