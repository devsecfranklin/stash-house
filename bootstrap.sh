#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 02/25/2022 Maintainer script
# v0.2 09/24/2022 Update this script
# v0.3 10/19/2022 Add tool functions
# v0.4 11/10/2022 Add automake check
# v0.5 11/15/2022 Handle Docker container builds

set -euo pipefail
IFS=$'\n\t'

#Black        0;30     Dark Gray     1;30
#Red          0;31     Light Red     1;31
#Green        0;32     Light Green   1;32
#Brown/Orange 0;33     Yellow        1;33
#Blue         0;34     Light Blue    1;34
#Purple       0;35     Light Purple  1;35
#Cyan         0;36     Light Cyan    1;36
#Light Gray   0;37     White         1;37

RED='\033[0;31m'
LRED='\033[1;31m'
LGREEN='\033[1;32m'
LBLUE='\033[1;34m'
CYAN='\033[0;36m'
LPURP='\033[1;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

MY_OS="unknown"
CONTAINER=false
DOCUMENTATION=false

# Check if we are inside a docker container
function check_docker() {
  if [ -f /.dockerenv ]; then
    log_warn "Building in container"
    CONTAINER=true
  fi
}

function check_installed() {
  if ! command -v ${1} &> /dev/null
  then
    log_warn "${1} could not be found"
    exit
  fi
}

function main() {
  check_docker
}

main "$@"
