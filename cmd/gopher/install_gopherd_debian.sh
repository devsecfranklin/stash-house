#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

SYSTEMDFILE="/etc/systemd/system/gopher-ctf.service"

MY_IP=$(host myip.opendns.com resolver1.opendns.com | grep address | cut -d' ' -f5)
echo "My IP address: ${MY_IP}"

function install_python() {
  install gopher
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
}

function main() {
  if [ -f "./gopher_server.go" ]; then  go build; fi # builds the "gopher" binary.
  cp gopher "${HOME}/.local.bin"
  sudo touch "${SYSTEMDFILE}" && sudo chown $(whoami) "${SYSTEMDFILE}"
  cat > "${SYSTEMDFILE}" <<EOF

  [Unit]
Description=Gopher CTF Challenge Server
# This ensures the service starts after the network is ready
After=network.target

[Service]
# It's good practice to run services as a non-root user.
# You can create a dedicated user with 'sudo adduser ctfuser'
# Or, for simplicity, use a non-root user that already exists, like 'ubuntu' or 'debian'.
# If you must run as root (e.g., to use port 70), you can comment this line out.
User=nobody

# The command to execute. We use 'screen' to launch our binary.
# -S gopher-ctf : Creates a screen session named "gopher-ctf"
# -d -m         : Starts the session in a detached state
# /usr/local/bin/gopher_server : The full path to our compiled binary
ExecStart=/usr/bin/screen -S gopher-ctf -d -m /usr/local/bin/gopher_server

# The command to gracefully stop the service.
# This command sends the "quit" command to the screen session, which terminates it.
ExecStop=/usr/bin/screen -S gopher-ctf -X quit

# Automatically restart the service if it crashes.
Restart=always

[Install]
# This tells systemd to start the service at boot time for multi-user mode.
WantedBy=multi-user.target
EOF
  sudo chown root "${SYSTEMDFILE}"
  sudo systemctl daemon-reload
  sudo systemctl enable gopher-ctf.service
  sudo systemctl start gopher-ctf.service
  sudo systemctl status gopher-ctf.service
}