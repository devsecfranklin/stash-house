#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 06/11/2025 GoLang Project Maintainer script

function main() {
  go get github.com/mattn/go-sqlite3
  go install github.com/kisielk/errcheck@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  test/test.sh
}

main "$@"
