#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 07/15/2026 initial version

if [ ! -d "-/.venv" ]; then 
  echo "run venv"
  python3 -m venv .venv
fi
source .venv/bin/activate && python3 -m pip install -r requirements.txt --break-system-packages

sudo apt-get install --reinstall texlive-base -y
sudo texhash
mkdir ~/texmf
tlmgr init-usertree
#tlmgr update --self # update
# Cross release updates are only supported with
# update-tlmgr-latest.sh --update
wget https://mirror.ctan.org/systems/texlive/tlnet/update-tlmgr-latest.sh && chmod 755 update-tlmgr-latest.sh
sh update-tlmgr-latest.sh -- --upgrade
tlmgr install montserrat
tlmgr install collection-fontsrecommended
tlmgr search --file montserrat.sty
