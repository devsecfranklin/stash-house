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

# Check if we are inside a container
CONTAINER=false
function check_container() {
  if [ -f /.dockerenv ]; then
    log_warn "Containerized build environment"
    CONTAINER=true
  else
    log_info "NOT a containerized build environment."
  fi
}

function check_installed() {
  if command -v "$1" &>/dev/null; then
    log_success "Found command: ${1}"
    return 0
  else
    log_error "Command not found: ${1}"
    return 1
  fi
}

PRIV_CMD="sudo"
function install_debian() {
  # Container package installs will fail unless you do an initial update, the upgrade is optional
  if [ "${CONTAINER}" = true ]; then
    log_info "Upgrading container packages"
    sudo apt-get update && apt-get upgrade -y
    sudo apt-get autoremove -y
  fi

  for i in "${DEB_PKG[@]}"; do
    if [ $(dpkg-query -W -f='${Status}' ${i} 2>/dev/null | grep -c "ok installed") -eq 0 ]; then
      log_warn "Installing ${i} since it is not found."
      # If we are in a container there is no sudo in Debian
      if [ "${CONTAINER}" = true ]; then
        $PRIV_CMD apt-get --yes install "${i}"
        $PRIV_CMD apt-get autoremove -y
      else
        $PRIV_CMD apt-get install "${i}" -y
        $PRIV_CMD apt-get autoremove -y
      fi
    else
      log_info "found package: ${i}"
    fi
  done

  if ! check_installed dircolors && [ ! -d "${HOME}/.dircolors" ]; then
    dircolors -p >~/.dircolors
    log_warn "Updating the dircolors configuration."
  fi
}
