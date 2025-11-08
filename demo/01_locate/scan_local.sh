#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#

LRED=$(tput setaf 1)

if [ -f "../../bin/common.sh" ]; then
  source "../../bin/common.sh"
else
  echo -e "${LRED}can not find common.sh${NC}"
  exit 1
fi

function scan_aws_creds() {
  log_info "scan for AWS credentials"

  if [ -d "${HOME}/.aws" ]; then
    while IFS=' = ' read -r key value; do
      if [[ $key == \[*] ]]; then
        section=$key
      elif [[ $value ]] && [[ $section == '[default]' ]]; then
        if [[ $key == 'aws_access_key_id' ]]; then
          MY_AWS_ACCESS_KEY_ID=$value
        elif [[ $key == 'aws_secret_access_key' ]]; then
          MY_AWS_SECRET_ACCESS_KEY=$value
        fi
      fi
    done <"${HOME}/.aws/credentials"

    if [ ! -z "${MY_AWS_ACCESS_KEY_ID}" ] || [ ! -z "${MY_AWS_SECRET_ACCESS_KEY}" ]; then
      # echo -e "\n"
      log_warn "FOUND AWS KEY: ${MY_AWS_ACCESS_KEY_ID}"
      log_warn "FOUND AWS SECRET: ${MY_AWS_SECRET_ACCESS_KEY}"
    else
      log_success "No local copies of AWS credentials found in standard location. Yay."
    fi
  fi
}

function scan_gcp_creds() {
  log_info "scan for GCP creds"

  if [ -d "${HOME}/.config/gcloud" ]; then
    log_info "found GCP folder"
  fi
}

function scan_gnome_keyring() {
  log_header "Scan Gnome keyring"
  # GNOME Keyring is a software application designed to store security
  # credentials such as usernames, passwords, and keys. The sensitive
  # data is encrypted and stored in a keyring file in the user’s home
  # directory. It can be found through the following command:
  locate login.keyring
  locate user.keystore
}

function scan_other() {
  log_info "Scan for the string PASSWORD"
  grep -rnw / -ie "PASSWORD\|PASSWD"
  find . -type f -exec grep -i -I "PASSWORD\|PASSWD" {} /dev/null \;

  log_info "Search for history files..."
  find / -name "*_history" -xdev 2>/dev/null

  log_info "Search for recent files..."
  find / -mmin -30 -xdev 2>/dev/null

  # strings /dev/mem -n10 | grep -ie “PASSWORD|PASSWD” –color=always
}

function crack_hash() {
  # John the Ripper can then be used to extract and crack the hashes and reveal the actual password:
  #/usr/share/john/keyring2john.py login.keyring > hashes.txt
  #/usr/share/john/keystore2john.py user.keystore
  #john –wordlist=/usr/share/wordlists/rockyou.txt hashes.txt

  #MimiPenguin and the post/linux/gather/gnome_keyring_dump Metasploit module can also be used to perform this task.
  pass
}

function main() {
  log_header "tool to show how local credential scanning works"
  # https://steflan-security.com/linux-privilege-escalation-credentials-harvesting/
  scan_aws_creds
  scan_gcp_creds # the JSON files for Google Cloud in ~/.config/gcloud
  # scan_gnome_keyring

}

main "$@"
