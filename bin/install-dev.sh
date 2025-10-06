#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 02/25/2025 initial version

DEB_PKG=(curl git gnupg2 make pass)


function main() {
  log_header "Installing Development Environment for the Stash House"
  check_python_version
}

main "$@"