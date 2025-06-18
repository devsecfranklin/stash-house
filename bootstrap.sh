#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 06/11/2025 GoLang Project Maintainer script

#set -euo pipefail

# The special shell variable IFS determines how Bash
# recognizes word boundaries while splitting a sequence of character strings.
#IFS=$'\n\t'

#Black        0;30     Dark Gray     1;30
#Red          0;31     Light Red     1;31
#Green        0;32     Light Green   1;32
#Brown/Orange 0;33     Yellow        1;33
#Blue         0;34     Light Blue    1;34
#Purple       0;35     Light Purple  1;35
#Cyan         0;36     Light Cyan    1;36
#Light Gray   0;37     White         1;37

#RED='\033[0;31m'
LRED='\033[1;31m'
LGREEN='\033[1;32m'
LBLUE='\033[1;34m'
CYAN='\033[0;36m'
LPURP='\033[1;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# --- Some config Variables ----------------------------------------
CONTAINER=false
GO_VERSION="$(go version | cut -d' ' -f3)"

# Check if we are inside a container
function check_container() {
  echo -e "\n${LPURP}# --- Check Container Status ------------------------------------------\n${NC}" | tee -a "${RAW_OUTPUT}"
  if [ -f /.dockerenv ]; then
    echo -e "${YELLOW}Containerized build environment...${NC}"
    CONTAINER=true
  else
    echo -e "${LBLUE}NOT a containerized build environment...${NC}"
  fi
}

function setup_golang() {
  echo "Go version: ${GO_VERSION}"

  if [ ! -f "go.mod" ]; then
    echo "initialize the go module"
    go mod init github.com/devsecfranklin/website
  fi

  echo "tidy up"
  go mod tidy

  go install github.com/mattn/go-sqlite3
  go install github.com/kisielk/errcheck@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  # test/test.sh
  if [ ! -f "${HOME}/go/bin/errcheck" ]; then go install github.com/kisielk/errcheck@latest; fi
  if [ ! -f "${HOME}/franklin/go/bin/golangci-lint" ]; then go install github.com/kisielk/errcheck@latest; fi
  go get internal/database
  go get internal/auth
  go get internal/cookies
}

function check_installed() {
  if ! command -v ${1} &>/dev/null; then
    echo "${1} could not be found"
    exit
  fi
}

function main() {
  check_container
  setup_golang
}

main "$@"
