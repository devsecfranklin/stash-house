#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2022-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: GPL-3.0-or-later

# ChangeLog:
#
# v0.1 02/25/2022 Maintainer script
# v0.2 09/24/2022 Update this script
# v0.3 10/19/2022 Add tool functions
# v0.4 11/10/2022 Add automake check
# v0.5 11/16/2022 Handle Docker container builds
# v0.6 07/13/2023 Add required_files and OpenBSD support
# v0.7 04/22/2024 More OpenBSD support
# v0.8 09/06/2024 Support GCP Linux
# v0.9 02/18/2025 Updates for Mac
# v1.0 02/26/2025 Optimize ssome functions using Gemini 2.0 Flash

#set -euo pipefail

# The special shell variable IFS determines how Bash
# recognizes word boundaries while splitting a sequence of character strings.
#IFS=$'\n\t'

#Black        0;30     Dark Gray     1;30
#Red          0;31     Light Red     1;31
#Green        0;32     Light Green   1;32
#Brown/Orange 0;33     Yellow        1;33
#Blue         0;34     Light Blue    1;34
#Purple       0;35     Light Purple  1;35
#Cyan         0;36     Light Cyan    1;36
#Light Gray   0;37     White         1;37

RED='\033[0;31m'
LRED='\033[1;31m'
LGREEN='\033[1;32m'
LBLUE='\033[1;34m'
CYAN='\033[0;36m'
LPURP='\033[1;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

MY_OS="unknown"
OS_RELEASE=""
CONTAINER=false
#DOCUMENTATION=false

# Check if we are inside a docker container
function check_docker() {
  if [ -f /.dockerenv ]; then
    echo -e "${CYAN}Containerized build environment...${NC}"
    CONTAINER=true
  else
    echo -e "${CYAN}NOT a containerized build environment...${NC}"
  fi
}

function detect_os() {
  # check for the /etc/os-release file
  if [ -f "/etc/os-release" ]; then
    OS_RELEASE=$(cat /etc/os-release | grep "^ID=" | cut -d"=" -f2)
  fi

  if [ -n "${OS_RELEASE}" ]; then
    echo -e "${CYAN}Found /etc/os-release file: ${OS_RELEASE}${NC}"
  fi

  # Check uname (Linux, OpenBSD, Darwin)
  MY_UNAME=$(uname)
  if [ -n "${OS_RELEASE}" ]; then
    echo -e "${CYAN}Found uname: ${MY_UNAME}${NC}"
  fi

  if [ "${MY_UNAME}" == "OpenBSD" ]; then
    echo -e "${CYAN}Detected OpenBSD${NC}"
    MY_OS="openbsd"
  elif [ "${MY_UNAME}" == "Darwin" ]; then
    echo -e "${CYAN}Detected MacOS${NC}"
    MY_OS="mac"
  elif [ -f "/etc/redhat-release" ]; then
    echo -e "${CYAN}Detected Red Hat/CentoOS/RHEL${NC}"
    MY_OS="rh"
  elif [ "$(grep -Ei 'debian|buntu|mint' /etc/*release)" ]; then
    echo -e "${CYAN}Detected Debian/Ubuntu/Mint${NC}"
    MY_OS="deb"
  elif grep -q Microsoft /proc/version; then
    echo -e "${CYAN}Detected Windows pretending to be Linux${NC}"
    MY_OS="win"
  else
    echo -e "${YELLOW}Unrecongnized architecture.${NC}"
    exit 1
  fi
}

function run_autopoint() {
  echo "Checking autopoint version..."
  ver=$(autopoint --version | awk '{print $NF; exit}')
  ap_maj=$(echo $ver | sed 's;\..*;;g')
  ap_min=$(echo $ver | sed -e 's;^[0-9]*\.;;g' -e 's;\..*$;;g')
  ap_teeny=$(echo $ver | sed -e 's;^[0-9]*\.[0-9]*\.;;g')
  echo "    $ver"

  case $ap_maj in
  0)
    if test $ap_min -lt 14; then
      echo "You must have gettext >= 0.14.0 but you seem to have $ver"
      exit 1
    fi
    ;;
  esac
  echo "Running autopoint..."
  autopoint --force || exit 1
}

function run_libtoolize() {
  echo "Checking libtoolize version..."
  libtoolize --version 2>&1 >/dev/null
  rc=$?
  if test $rc -ne 0; then
    echo "Could not determine the version of libtool on your machine"
    echo "libtool --version produced:"
    libtool --version
    exit 1
  fi
  lt_ver=$(libtoolize --version | awk '{print $NF; exit}')
  lt_maj=$(echo $lt_ver | sed 's;\..*;;g')
  lt_min=$(echo $lt_ver | sed -e 's;^[0-9]*\.;;g' -e 's;\..*$;;g')
  #lt_teeny=$(echo $lt_ver | sed -e 's;^[0-9]*\.[0-9]*\.;;g')
  echo "    $lt_ver"

  case $lt_maj in
  0)
    echo "You must have libtool >= 1.4.0 but you seem to have ${lt_ver}"
    exit 1
    ;;
  1)
    if test "${lt_min}" -lt 4; then
      echo "You must have libtool >= 1.4.0 but you seem to have ${lt_ver}"
      exit 1
    fi
    ;;
  2) ;;
  *)
    echo "You are running a newer libtool than gerbv has been tested with."
    echo "It will probably work, but this is a warning that it may not."
    ;;
  esac
  echo "Running libtoolize..."
  libtoolize --force --copy --automake || exit 1
}

function run_aclocal() {
  if [ "${MY_OS}" != "openbsd" ]; then
    echo -e "${LBLUE}Checking aclocal version...${NC}"
    acl_ver=$(aclocal --version | awk '{print $NF; exit}')
    echo "    $acl_ver"

    echo -e "${CYAN}Running aclocal...${NC}"
    #aclocal -I m4 $ACLOCAL_FLAGS || exit 1
    aclocal -Iaclocal/latex-m4 || exit 1
  else
    AUTOCONF_VERSION=2.71 AUTOMAKE_VERSION=1.16 aclocal -Iaclocal/latex-m4 || exit 1
  fi
  echo -e "${CYAN}.. done with aclocal.${NC}"
}

function run_autoheader() {
  echo "Checking autoheader version..."
  ah_ver=$(autoheader --version | awk '{print $NF; exit}')
  echo "    $ah_ver"

  echo "Running autoheader..."
  autoheader || exit 1
  echo "... done with autoheader."
}

function run_automake() {
  if [ "${MY_OS}" != "openbsd" ]; then
    echo "Checking automake version..."
    am_ver=$(automake --version | awk '{print $NF; exit}')
    echo "    $am_ver"

    echo "Running automake..."
    automake -a -c --add-missing || exit 1
    #automake --force --copy --add-missing || exit 1
  else
    AUTOCONF_VERSION=2.71 AUTOMAKE_VERSION=1.16 automake -a -c --add-missing || exit 1
  fi
  echo "... done with automake."
}

function run_autoconf() {
  if [ "${MY_OS}" != "openbsd" ]; then
    echo -e "${LGREEN}Checking autoconf version...${NC}"
    ac_ver=$(autoconf --version | awk '{print $NF; exit}')
    echo -e "${LGREEN}Autoconf version: $ac_ver${NC}"
    echo "Running autoconf..."
    autoreconf -i || exit 1
  else
    # this is for OpenBSD systems
    ac_ver="2.71"
    echo "Running autoconf..."
    AUTOCONF_VERSION=2.71 AUTOMAKE_VERSION=1.16 autoreconf -i || exit 1
  fi
  echo "... done with autoconf."
}

function check_installed() {
  if command -v "$1" &>/dev/null; then
    printf "${LPURP}Found command: %s${NC}\n" "$1"
    return 0
  else
    printf "${LRED}%s could not be found${NC}\n" "$1"
    return 1
  fi
}

function install_macos() {
  #declare -a Packages=("ac")
  declare -a Packages=("google-cloud-sdk" "git" "bash" "make" "automake" "gsed" "gawk" "direnv" "terraform" "libtool" "jq" "google-cloud-sdk" "coreutils")

  echo -e "${CYAN}Updating brew for MacOS (this may take a while...)${NC}"
  brew update
  #brew upgrade google-cloud-sdk # this is to avoid the error: ModuleNotFoundError: No module named 'imp'

  for i in "${Packages[@]}"; do
    if brew list "${i}" &>/dev/null; then
      echo -e "${LGREEN}${i} is already installed${NC}"
      brew upgrade "${i}"
    else
      brew install "${i}"
    fi
  done

  echo -e "${CYAN}Updating Google gcloud for MacOS (this may take a while...)${NC}"
  (yes || true) | "${HOME}/homebrew/bin/gcloud" components update

  if [ ! -f "./config.status" ]; then
    echo -e "${CYAN}Running libtool/autoconf/automake...${NC}"
    # glibtoolize
    aclocal -I config
    autoreconf -i
    automake -a -c --add-missing
  else
    echo -e "${CYAN}Your system is already configured. (Delete config.status to reconfigure)${NC}"
    ./config.status
  fi
  echo -e "${CYAN}HINT: now type \"./configure\"${NC}"

  # https://github.com/kreuzwerker/m1-terraform-provider-helper/blob/main/README.md
  brew install kreuzwerker/taps/m1-terraform-provider-helper
  brew tap hashicorp/tap
  brew install hashicorp/tap/terraform
  brew install yasm
  m1-terraform-provider-helper activate
  #m1-terraform-provider-helper install hashicorp/template -v v2.2.0 # DEPRECATED
  #terraform providers lock -platform=darwin_arm64
  #terraform providers lock -platform=linux_amd64

  echo -e "${CYAN}Running brew cleanup...${NC}"
  brew cleanup
}

function install_debian() {
  declare -a Packages=("libonig-dev" "sshpass" "shellcheck" "screen" "make" "gcc" "git" "automake" "libtool" "gawk" "nodejs" "npm" "apt-transport-https" "ca-certificates" "curl" "gnupg" "lsb-release") # "python3-pygit2" )

  # Container package installs will fail unless you do an initial update, the upgrade is optional
  if [ "${CONTAINER}" = true ]; then
    echo -e "${LBLUE}Upgrading container packages${NC}"
    apt-get update && apt-get upgrade -y
  fi

  for i in "${Packages[@]}"; do
    if [ $(dpkg-query -W -f='${Status}' ${i} 2>/dev/null | grep -c "ok installed") -eq 0 ]; then
      echo -e "${LBLUE}Installing ${i} since it is not found.${NC}"
      # If we are in a container there is no sudo in Debian
      if [ "${CONTAINER}" = true ]; then
        apt-get --yes install "${i}"
      else
        sudo apt-get install "${i}" -y
      fi
    fi
  done

  if check_installed dircolors && [ ! -f "${HOME}/.dircolors" ]; then
    dircolors -p >~/.dircolors
    echo -e "${LBLUE}Updating the dircolors configuration.${NC}"
  fi

}

function install_az_cli() {
  # Install az cli tool for Azure
  #curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
  # sudo apt-get update
  # sudo apt-get install apt-transport-https ca-certificates curl gnupg lsb-release
  # sudo mkdir -p /etc/apt/keyrings
  # curl -sLS https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor | sudo tee /etc/apt/keyrings/microsoft.gpg > /dev/null
  # sudo chmod go+r /etc/apt/keyrings/microsoft.gpg
  # AZ_DIST=$(lsb_release -cs)
  # echo "Types: deb
  # URIs: https://packages.microsoft.com/repos/azure-cli/
  # Suites: ${AZ_DIST}
  # Components: main
  # Architectures: $(dpkg --print-architecture)
  # Signed-by: /etc/apt/keyrings/microsoft.gpg" | sudo tee /etc/apt/sources.list.d/azure-cli.sources
  # sudo apt-get update
  # sudo apt-get install azure-cli
  echo "" #pass
}

function debian() {
  # sudo apt install gnuplot gawk libtool psutils make autopoint
  # run_autopoint
  run_aclocal
  autoreconf -i
  run_automake
  ./configure
  #./config.status
}

function redhat() {
  if [ ! -f "./config.status" ]; then
    mkdir -p aclocal # Create aclocal if needed
    run_aclocal
    autoreconf -i
    run_automake
    ./configure
  else
    ./config.status
  fi
}

function install_redhat() {
  echo -e "${CYAN}RedHat 8 setup${NC}"
  dnf upgrade -y
  yum -y --disableplugin=subscription-manager update
  dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm

  declare -a Packages=("make" "automake" "autoconf" "libtool" "texlive")
  for i in "${Packages[@]}"; do
    dnf install -y "${i}" --skip-broken
  done
}

function required_files() {
  echo "Check for presence of required GNU autotools files"

  local required_files=("AUTHORS" "ChangeLog" "NEWS")

  for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
      printf "${LGREEN}Creating required file %s since it is not found.${NC}\n" "$file"
      ln -sf README.md "$file" # Use -sf to force and be silent
    else
      printf "${LBLUE}Found required file %s.${NC}\n" "$file"
    fi
  done

  mkdir -p config/m4 # Create config/m4 if needed
}

function main() {
  check_docker
  detect_os
  # check_installed doxygen
  required_files

  if [ ! -d "aclocal" ]; then mkdir aclocal; fi
  if [ ! -d "config/m4" ]; then mkdir -p config/m4; fi

  if [ "${MY_OS}" == "mac" ]; then
    check_installed brew
    install_macos
  fi

  if [ "${MY_OS}" == "rh" ]; then
    install_redhat
  fi

  if [ "${MY_OS}" == "deb" ]; then
    install_debian
  fi

  if [ ! -f "Makefile.in" ] && [ -f "./config.status" ]; then
    rm config.status # if Makefile.in is missing, then erase stale config.status
  fi

  if [ ! -f "./config.status" ]; then
    echo -e "${YELLOW}no config.status${NC}"
    # libtoolize
    if [ ! -d "aclocal" ]; then mkdir aclocal; fi
    #aclocal -I config
    run_aclocal
    if [ "${MY_OS}" == "openbsd" ]; then
      AUTOCONF_VERSION=2.71 AUTOMAKE_VERSION=1.16 autoreconf -i || exit 1
    else
      autoreconf -i
    fi
    #automake -a -c --add-missing
    run_automake
    ./configure
  else
    ./config.status
  fi
}

main "$@"
