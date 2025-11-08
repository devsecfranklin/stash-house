#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 11/07/2025 initial version

LRED=$(tput setaf 1)
MY_KEY="9FBB57D0A73805E48BCC1F67F518C4034CAC31C2"

function key_select() {
  log_header "Show your existing keys:"
  gpg --list-secret-keys --keyid-format=long

  KEYS=$(gpg --list-secret-keys --keyid-format=long | tail -n +3 | grep -v sec | grep -v uid | grep -v ssb | grep .)
  KEY_LIST=($KEYS)
  for i in "${KEY_LIST[@]}"; do
    log_info "Found key: ${i}"
  done
}

function key_upload() {
  log_header "Upload your key: ${MY_KEY}"

  gpg --keyserver hkp://keyserver.ubuntu.com --send-key "${MY_KEY}"
  gpg --keyserver hkp://keyserver.openpgp.org --send-key "${MY_KEY}"
  gpg --keyserver u --send-key "${MY_KEY}"
  gpg --keyserver hkp://keys.openpgp.org --send-key "${MY_KEY}"
}

# https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key
# log_header "Generate a new key"
# gpg --full-generate-key

# function key_delete(){
# Revoke the key - by creating revocation certificate
# delete }

function main() {

  if [ -f "../../bin/common.sh" ]; then
    source "../../bin/common.sh"
  else
    echo -e "${LRED}can not find common.sh${NC}"
    exit 1
  fi
  key_select
  key_upload
}

main "$@"
