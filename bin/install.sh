#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 02/25/2025 initial version

# SCRIPT_DIR="${0%/*}"
if [ -f "./common.sh" ]; then
  source "./common.sh"
else
  echo "can not find common.sh"
  exit 1
fi

function main() {
  log_header "Installing Stash House"
  check_if_root
}

main "$@"
