#!/bin/bash

MY_IP=$(host myip.opendns.com resolver1.opendns.com | grep address | cut -d' ' -f5)

apt install gopher -y 

if [ ! -d "/var/gopher" ]; then mkdir /var/gopher && chmod 755 /var/gopher; fi

cp files/gophermap /var/gopher

if [ ! -d "/etc/pygopherd" ]; then mkdir /etc/pygopherd && chmod 755 /var/gopher; fi

cp files/pygopherd.conf /etc/pygopherd/pygopherd.conf


