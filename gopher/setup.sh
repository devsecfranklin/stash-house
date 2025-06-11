#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

MY_IP=$(host myip.opendns.com resolver1.opendns.com | grep address | cut -d' ' -f5)
echo "My IP address: ${MY_IP}"

# install gopher
# sudo pip install -rrequirements.txt
sudo apt install gopher -y 

# create gopher folder
if [ ! -d "/var/gopher" ]; then
  sudo mkdir /var/gopher
  sudo chmod 755 /var/gopher
fi

# files
cp files/gophermap files/*.txt /var/gopher

# pygopherd
if [ ! -d "/etc/pygopherd" ]; then
  sudo mkdir /etc/pygopherd
  sudo chmod 755 /var/gopher
fi
sudo cp files/pygopherd.conf /etc/pygopherd/pygopherd.conf
