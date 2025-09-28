#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 02/25/2022 Maintainer script
# v0.2 09/24/2022 Update this script
# v0.3 10/19/2022 Add tool functions
# v0.4 11/10/2022 Add automake check
# v0.5 11/15/2022 Handle Docker container builds


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
