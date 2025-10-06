#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:
#

# --- Script Configuration ---
# Exit immediately if a command exits with a non-zero status.
set -e
# Treat unset variables as an error.
set -u
# The return value of a pipeline is the status of the last command to exit
# with a non-zero status, or zero if no command failed.
set -o pipefail

# --- Color and Logging Functions ---
# Using tput for compatibility and to check if the terminal supports color.
if tput setaf 1 &> /dev/null; then
    RED=$(tput setaf 1)
    GREEN=$(tput setaf 2)
    YELLOW=$(tput setaf 3)
    CYAN=$(tput setaf 6)
    LPURP='\033[1;35m'
    BOLD=$(tput bold)
    NC=$(tput sgr0) # No Color
else
    RED=""
    GREEN=""
    YELLOW=""
    CYAN=""
    BOLD=""
    NC=""
fi

# Centralized logging functions for consistent output.
log_header() { printf "\n${LPURP}# --- %s ${NC}\n" "$1"; }
log_info() { echo -e "${CYAN}==>${NC} ${BOLD}$1${NC}"; }
log_success() { echo -e "${GREEN}==>${NC} ${BOLD}$1${NC}"; }
log_warn() { echo -e "${YELLOW}WARN:${NC} $1"; }
log_error() { >&2 echo -e "${RED}ERROR:${NC} $1"; } # Errors to stderr

function check_if_root {
  if [[ $(id -u) -eq 0 ]]; then
    log_warn "You are the root user."
  else
    log_success "You are NOT the root user."
  fi
}