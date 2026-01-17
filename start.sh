#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#
# v0.1 | 29 Sept 2023 | franklin <smoooth.y62wj@passmail.net>
# v0.2 | 12/27/2025 | update for common.sh
# v0.3 | 01/16/2026 | 

set -o nounset  # Treat unset variables as an error
#set -euo pipefail

# The special shell variable IFS determines how Bash
# recognizes word boundaries while splitting a sequence of character strings.
#IFS=$'\n\t'

LRED='\033[0;31m'
NC='\033[0m' # No Color
WORK_DIR="${PWD}"


function main() {
  echo -e "\n" && figlet -f ${WORK_DIR}/fonts/pagga resume && echo -e "\n"
  if [ -f "${WORK_DIR}/bin/common.sh" ]; then
    source "${WORK_DIR}/bin/common.sh"
  else
    echo -e "${LRED}can not find ${WORK_DIR}/bin/common.sh.${NC}"
    exit 1
  fi
  log_info "successfully sourced ${WORK_DIR}/bin/common.sh" && echo -e "\n"

  log_header "Starting the server."

  log_info "cd to ${WORK_DIR}"
  pushd "${WORK_DIR}" >> /dev/null || exit 1       
  log_info "start screen session"
  screen -mdS website go run ${WORK_DIR}/cmd/www/main.go
  popd >> /dev/null || exit 1
  log_success "Open in browser: http://localhost:9091"
}

main "$@"
