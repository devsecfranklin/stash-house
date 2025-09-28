#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#

# SCRIPT_DIR="${0%/*}"
if [ -f "bin/common.sh" ]; then
  source "bin/common.sh"
else
    log_error "can not find common.sh"
fi

function scan_aws_creds() {
  log_info "scan for AWS credentials"

  while IFS=' = ' read key value
  do
      if [[ $key == \[*] ]]; then
          section=$key
      elif [[ $value ]] && [[ $section == '[default]' ]]; then
          if [[ $key == 'aws_access_key_id' ]]; then
              AWS_ACCESS_KEY_ID=$value
          elif [[ $key == 'aws_secret_access_key' ]]; then
              AWS_SECRET_ACCESS_KEY=$value
          fi
      fi
  done < "${HOME}/.aws/credentials"

  # echo -e "\n"
  log_warn "FOUND AWS KEY: $AWS_ACCESS_KEY_ID"
  log_warn "FOUND AWS SECRET: $AWS_SECRET_ACCESS_KEY"
}

function scan_gnome_keyring() {
  # GNOME Keyring is a software application designed to store security
  # credentials such as usernames, passwords, and keys. The sensitive
  # data is encrypted and stored in a keyring file in the user’s home
  # directory. It can be found through the following command:
  locate login.keyring
  locate user.keystore
}

function main(){
# https://steflan-security.com/linux-privilege-escalation-credentials-harvesting/
scan_aws_creds
scan_gnome_keyring

exit 1

grep -rnw / -ie "PASSWORD\|PASSWD"

find . -type f -exec grep -i -I "PASSWORD\|PASSWD" {} /dev/null \;

# history files
find / -name *_history -xdev

# recent files
find / -mmin -30 -xdev 2>/dev/null

# the JSON files for Google CLoud in ~/.config/gcloud

# strings /dev/mem -n10 | grep -ie “PASSWORD|PASSWD” –color=always



# John the Ripper can then be used to extract and crack the hashes and reveal the actual password:
#/usr/share/john/keyring2john.py login.keyring > hashes.txt
#/usr/share/john/keystore2john.py user.keystore
#john –wordlist=/usr/share/wordlists/rockyou.txt hashes.txt

#MimiPenguin and the post/linux/gather/gnome_keyring_dump Metasploit module can also be used to perform this task.
}

main "$@"
