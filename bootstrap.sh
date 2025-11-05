#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:

DEB_PKG=(libpcsclite-dev texlive-pictures texlive-latex-extra)
LRED=$(tput setaf 1)

function setup_foks() {
  log_header "setup foks"
  /usr/local/go/bin/go install github.com/foks-proj/go-foks/client/foks@latest
}

if [ -f "./bin/common.sh" ]; then
  source "./bin/common.sh"
  log_success "Source common routines: bin/common.sh"
else
  echo -e "${LRED}can not find common.sh, run this script from the bin/ directory.${NC}"
  exit 1
fi

function create_files() {
  log_header "Create files for gnu compiler"
  FILES=(AUTHORS ChangeLog INSTALL NEWS)
  for i in "${FILES[@]}"; do
    touch "${i}"
  done
} 

function main() {
  install_debian
  create_files
  setup_foks
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

log_success "Complete!"
 }


main "$@"
