#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:

DEB_PKG=(libpcsclite-dev texlive-pictures texlive-latex-extra)
LRED=$(tput setaf 1)

GOBIN="${GOROOT}/bin"

function setup_foks() {
  log_header "setup foks"
  ${GOBIN}/go install github.com/foks-proj/go-foks/client/foks@latest
}

if [ -f "./bin/common.sh" ]; then
  source "./bin/common.sh"
  log_success "Source common routines: bin/common.sh"
else
  echo -e "${LRED}can not find common.sh, run this script from the bin/ directory.${NC}"
  exit 1
fi

function create_files() {
  log_header "Setup:: GNU Compiler/GNU Autotools"
  FILES=(README AUTHORS ChangeLog INSTALL NEWS)
  for i in "${FILES[@]}"; do
    touch "${i}"
  done
}

function main() {
  install_debian
  create_files
  if [ ! -d "aclocal" ]; then
    log_info "create aclocal dir"
    mkdir -p aclocal
  fi

  if [ ! -d "config/m4" ]; then
    log_info "create congif/m4 dir"
    mkdir -p config/m4
  fi

  if [ ! -f "Makefile.in" ] && [ -f "./config.status" ]; then
    log_info "remove config.status"
    rm config.status # if Makefile.in is missing, then erase stale config.status
  fi

  log_info "Create the Makefiles"
  aclocal -Iaclocal/latex-m4 || exit 1
  autoreconf -i
  automake -a -c --add-missing || exit 1
  if [ -f "./configure" ]; then
    log_info "Running configure..."
    ./configure
  fi

  
  # ./config.status

  log_header "Setup:: Go"
  if [ ! -f "go.mod" ]; then
    log_warn "creating go.mod"
    ${GOBIN}/go mod init github.com/devsecfranklin/stash-house
  else
    log_info "Found go.mod. Nice."
  fi

  log_info "Update all dependencies"
  ${GOBIN}/go get -u ./...

  log_info "remove unused deps"
  ${GOBIN}/go mod tidy

  log_info "viper module for config files"
  ${GOBIN}/go get github.com/spf13/viper

  # gorm is to connect to amriadb
  ${GOBIN}/go get -u gorm.io/gorm
  ${GOBIN}/go get -u gorm.io/driver/mysql 

  # setup_foks
  log_success "Complete!"
}

main "$@"
