#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# --- Color and Logging Functions ---
# Use tput for compatibility and to check if the terminal supports color.
setup_colors() {
    if tput setaf 1 &>/dev/null; then
        RED=$(tput setaf 1)
        GREEN=$(tput setaf 2)
        YELLOW=$(tput setaf 3)
        BLUE=$(tput setaf 4)
        PURPLE=$(tput setaf 5)
        CYAN=$(tput setaf 6)
        BOLD=$(tput bold)
        NC=$(tput sgr0) # No Color
    else
        RED=""
        GREEN=""
        YELLOW=""
        BLUE=""
        PURPLE=""
        CYAN=""
        BOLD=""
        NC=""
    fi
}

# Centralized logging functions for consistent output.
log_info() { echo -e "${CYAN}==>${NC} ${BOLD}$1${NC}"; }
log_success() { echo -e "${GREEN}==>${NC} ${BOLD}$1${NC}"; }
log_warn() { >&2 echo -e "${YELLOW}WARN:${NC} $1"; }
log_error() { >&2 echo -e "${RED}ERROR:${NC} $1"; }

for i in $(ls -d */); do
  pushd $i || fail "unable to find dir: $i" 1
  log_info "updating $i"
  /usr/local/go/bin/go mod init "bitsmasher.net/${i}"
  /usr/local/go/bin/go mod tidy
  /usr/local/go/bin/go mod verify
  popd || exit 1
done

