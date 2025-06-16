#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 06/11/2025 GoLang Project Maintainer script

function setup_golang() {
  if [ ! -f "${HOME}/go/bin/errcheck" ]; then go install github.com/kisielk/errcheck@latest; fi
  if [ ! -f "${HOME}/franklin/go/bin/golangci-lint" ]; then go install github.com/kisielk/errcheck@latest; fi
}

function main() {
  
  if [ ! -f "go.mod" ]; then
    go mod init github.com/devsecfranklin/website
    go mod tidy
  fi

  go install github.com/mattn/go-sqlite3
  go install github.com/kisielk/errcheck@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  # test/test.sh
}

main "$@"
