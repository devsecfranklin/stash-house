#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2022-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

set -euo pipefail

# The special shell variable IFS determines how Bash
# recognizes word boundaries while splitting a sequence of character strings.
IFS=$'\n\t'


# --- Some config Variables ----------------------------------------
TOP_DIR="$(pwd)"
MY_DATE=$(date '+%Y-%m-%d-%H')
RAW_OUTPUT="/tmp/bootstrap_website_${MY_DATE}.log"


function install_check_golang_templates(){
  git clone git@github.com:gosom/check-golang-templates.git /tmp/check-golang-templates 
  pushd /tmp/check-golang-templates || exit 1
  go mod download
  go install
  popd /tmp/check-golang-template || exit 1
}

function validate_templates() {
  log_header "# --- Linting ----------------------------------------------"
  golangci-lint run

  HOSTS=(games gopher www lab.bitsmasher.net)
  for MY_HOST in "${HOSTS[@]}"; do
    log_info "# --- Run Errcheck on ${MY_HOST} ----------------------------------------------"
    # cd "${TOP_DIR}" || exit 1
    # "${HOME}/go/bin/errcheck" ./...
    "${HOME}/go/bin/errcheck" "cmd/${MY_HOST}/"

    log_info "# --- Validate Templates ${MY_HOST} ----------------------------------------------"
    check-golang-templates -folder "template/${MY_HOST}"
  done
}

function execute_tests() {
  log_header "# --- Test www ----------------------------------------------"
  # cd "${WWW_DIR}" || exit 1
  # if [ -f "cmd/www/main.go" ]; then
  #   "${HOME}/go/bin/errcheck" ./...
  # fi

  # npx linthtml sites/**/*.html
  go test -v unit/main_test.go
  # go test -v unit/gopher_server_test.go
}

function main() {

  # SCRIPT_DIR="${0%/*}"
  if [ -f "../bin/common.sh" ]; then
    source "../bin/common.sh"
  else
    echo "can not find common.sh"
    exit 1
  fi

  log_info "Install testify/assert"
  go get github.com/stretchr/testify/assert
  log_info "Install testify/require"
  go get github.com/stretchr/testify/require

  setup_figlet

  # install_check_golang_templates

  if [ ! -f "go.mod" ]; then
    go mod init github.com/devsecfranklin/website
    go mod tidy
  fi


}

main "$@"
