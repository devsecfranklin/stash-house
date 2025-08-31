#!/bin/bash

# v0.1 11/18/2022 first version

# Author:      Franklin D. <devsecfranklin@duck.com>

#set -eu

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

CONTAINER=false

# Check if we are inside a docker container
function check_docker() {
  if [ -f /.dockerenv ]; then
    echo -e "${LGREEN}Building in container...${NC}"
    CONTAINER=true
  else
    echo -e "${LGREEN}NOT Building in container...${NC}"
  fi
}

# this will check if package is installed.
function check_installed() {
  if ! command -v ${1} &> /dev/null
  then
    echo -e "${LRED}${1} could not be found${NC}"
    exit
  else
    echo -e "${CYAN}Found package ${1} installed.${NC}"
  fi
}

# These are known packages we need on Debian
function install_debian() {
  declare -a Packages=( "gnupg" "pass" "git")
  
  # Container package installs will fail unless you do an initial update, the upgrade is optional
  if [ "${CONTAINER}" = true ]; then
    apt-get update && apt-get upgrade -y
  fi
  
  for i in ${Packages[@]};
  do
    PKG_OK=$(dpkg-query -W --showformat='${Status}\n' ${i}|grep "install ok installed") &> /dev/null
    # echo -e "${LBLUE}Checking for ${i}: ${PKG_OK}${NC}"
    if [ "" = "${PKG_OK}" ]; then
      echo -e "${CYAN}Installing ${i} since it is not found.${NC}"
      
      # If we are in a container there is no sudo in Debian
      if [ "${CONTAINER}" = true ]; then
        apt-get --yes install ${i}
      else
        sudo apt-get install ${i} -y
      fi
    fi
  done
}

function make_directories() {
  DIRS=("~/.password-store" "~/.password-store/documents")
  for i in ${DIRS[@]};
  do
    if [ ! -d "${i}" ]; then mkdir -p ${i}; fi
  done
}

# Rustlang app
function build_application() {
  cargo run stash
  target/debug/stash
}

function detect_os() {
  if [ "$(uname)" == "Darwin" ]
  then
    echo -e "${CYAN}Detected MacOS${NC}"
    MY_OS="mac"
  elif [ -f "/etc/redhat-release" ]
  then
    echo -e "${CYAN}Detected Red Hat/CentoOS/RHEL${NC}"
    MY_OS="rh"
  elif [ "$(grep -Ei 'debian|buntu|mint' /etc/*release)" ]
  then
    echo -e "${CYAN}Detected Debian/Ubuntu/Mint${NC}"
    MY_OS="deb"
  elif grep -q Microsoft /proc/version
  then
    echo -e "${CYAN}Detected Windows pretending to be Linux${NC}"
    MY_OS="win"
  else
    echo -e "${YELLOW}Unrecongnized architecture.${NC}"
    exit 1
  fi
}

function install_macos() {
  echo -e "${CYAN}Updating brew for MacOS (this may take a while...)${NC}"
  check_installed brew
  brew cleanup
  brew upgrade
}

function install_debian() {
  check_installed gpg
  check_installed pass
}

function check_rcs() {
  if [ -d "${HOME}/.password-store/.git" ]; then
    REPO=$(cat ${HOME}/.password-store/.git/config | grep "url = " | sed 's:.*=::' | xargs)
    echo -e "${LGREEN}Found repo ${CYAN}${REPO}${NC}"
  else
    echo -e "${RED}No Github Repo found${NC}"
    exit 1
  fi
}

function main() {
  check_docker
  detect_os
  
  if [ "${MY_OS}" == "mac" ]
  then
    install_macos
  elif [ "${MY_OS}" == "deb" ]
  then
    install_debian
  else
    echo -e "${RED}Unsupported OS, open a ticket on Git repo for script update${NC}"
    exit 1
  fi
  
  check_rcs

  # RustLang version
  #check_installed rustc
  #build_application
}

main
