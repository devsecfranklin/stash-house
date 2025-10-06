#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 02/25/2025 initial version


DEB_PKG=(pass gnupg)
LOCAL_STASH="${HOME}/.stash-house"
#SCRIPT_DIR="${0%/*}" # echo "$SCRIPT_DIR"


if [ -f "./common.sh" ]; then
  source "./common.sh"
else
  echo "can not find common.sh"
  exit 1
fi

function install_local() {
  
  log_header "Install the local passwd storage folder"

  if [ ! -d "${LOCAL_STASH}" ]; then
    mkdir -p "${LOCAL_STASH}"
    log_info "${LOCAL_STASH} directory was created."
  else 
    log_info "${LOCAL_STASH} directory already exists."
  fi
  
  if [ ! -f "${LOCAL_STASH}/.gitattributes" ]; then
    echo "*.gpg diff=gpg" > "${LOCAL_STASH}/.gitattributes"
    log_info "created file: ${LOCAL_STASH}/.gitattributes"
  fi

  check_installed pass
  check_installed keyringer
}


function main() {
  log_header "Installing Stash House"
  check_if_root
  check_container
  install_local
  # install_debian
}

main "$@"
