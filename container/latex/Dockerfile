# syntax=docker/dockerfile:1

<<<<<<< HEAD
FROM python:3.13.5-bullseye
=======
FROM python:3.13.6-bullseye
>>>>>>> main

#####################
# Add some packages #
#####################
ENV DEBIAN_FRONTEND noninteractive
RUN \
    apt-get update; \
    apt-get install -y apt-utils; \
    apt-get install -y make libgraphviz-dev lacheck
