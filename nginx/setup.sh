#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 02/25/2022 Maintainer script
#
sudo apt install -y nginx

sudo /usr/sbin/nginx -t


function test_nostr() {
  curl -k -L http://bitsmasher.net/.well-known/nostr.json
  ls -al /usr/share/nginx/html/.well-known/
}
