#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2022-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: GPL-3.0-or-later

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

RED='\033[0;31m'
LRED='\033[1;31m'
LGREEN='\033[1;32m'
LBLUE='\033[1;34m'
CYAN='\033[0;36m'
LPURP='\033[1;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# --- Some config Variables ----------------------------------------
TOP_DIR="$(pwd)"
GAMES_DIR="${TOP_DIR}/sites/games"
WWW_DIR="${TOP_DIR}/sites/www"
MY_DATE=$(date '+%Y-%m-%d-%H')
RAW_OUTPUT="/tmp/bootstrap_website_${MY_DATE}.log"

function check_installed() {
  if command -v "$1" &>/dev/null; then
    printf "${LPURP}Found command: %s${NC}\n" "$1"
    return 0
  else
    printf "${LRED}%s could not be found${NC}\n" "$1"
    return 1
  fi
}

function setup_figlet() {
  if ! check_installed figlet; then sudo apt update && sudo apt install figlet -y; fi
  echo -e "${CYAN}Figlet font setup${NC}"
  if [ ! -d "/tmp/figlet-fonts" ]; then cd /tmp && git clone https://github.com/xero/figlet-fonts.git; fi
  if [ ! -d "/usr/share/figlet/fonts" ]; then sudo mv /tmp/figlet-fonts /usr/share/figlet/fonts; fi
  figlet -f /usr/share/figlet/fonts/smmono9 golang test
}

function setup_golang() {
  if [ ! -f "${HOME}/go/bin/errcheck" ]; then go install github.com/kisielk/errcheck@latest; fi
  if [ ! -f "${HOME}/franklin/go/bin/golangci-lint" ]; then go install github.com/kisielk/errcheck@latest; fi
}

function main() {
  setup_figlet
  # The git repo has a .git dir, while a submodule has a .git file
  if [ ! -f "./.git" ]; then
    echo -e "${RED}ERROR: ${YELLOW}Run script from top level of your Git repo${NC}"
    exit 1
  fi

  echo -e "${LPURP}# --- Test games ----------------------------------------------\n${NC}" | tee -a "${RAW_OUTPUT}"
  cd "${GAMES_DIR}" || exit 1
  if [ -f "cmd/www/main.go" ]; then
    go mod tidy
    "${HOME}/go/bin/errcheck" ./...
    golangci-lint run
  fi

  # echo -e "${LPURP}# --- Test www ----------------------------------------------\n${NC}" | tee -a "${RAW_OUTPUT}"
  # cd "${WWW_DIR}" || exit 1
  # if [ -f "cmd/www/main.go" ]; then
  #   "${HOME}/go/bin/errcheck" ./...
  # fi
  npx linthtml sites/**/*.html
}

main "$@"
