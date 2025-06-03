#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

MY_IP=$(host myip.opendns.com resolver1.opendns.com | grep address | cut -d' ' -f5)
echo "${MY_IP}"

sudo apt install gopher -y 

if [ ! -d "/var/gopher" ]; then
  sudo mkdir /var/gopher
  sudo chmod 755 /var/gopher
fi

cp files/gophermap /var/gopher

if [ ! -d "/etc/pygopherd" ]; then
  sudo mkdir /etc/pygopherd
  sudo chmod 755 /var/gopher
fi

sudo cp files/pygopherd.conf /etc/pygopherd/pygopherd.conf