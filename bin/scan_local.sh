#!/usr/bin/env bash

# https://steflan-security.com/linux-privilege-escalation-credentials-harvesting/

#set -euo pipefail
#IFS=$'\n\t'


grep -rnw / -ie “PASSWORD\|PASSWD”

find . -type f -exec grep -i -I “PASSWORD\|PASSWD” {} /dev/null \;

# history files
find / -name *_history -xdev

# recent files
find / -mmin -30 -xdev 2>/dev/null

# the JSON files for Google CLoud in ~/.config/gcloud

# strings /dev/mem -n10 | grep -ie “PASSWORD|PASSWD” –color=always

# GNOME Keyring is a software application designed to store security
# credentials such as usernames, passwords, and keys. The sensitive
# data is encrypted and stored in a keyring file in the user’s home
# directory. It can be found through the following command:
locate login.keyring
locate user.keystore

# John the Ripper can then be used to extract and crack the hashes and reveal the actual password:
#/usr/share/john/keyring2john.py login.keyring > hashes.txt
#/usr/share/john/keystore2john.py user.keystore
#john –wordlist=/usr/share/wordlists/rockyou.txt hashes.txt

#MimiPenguin and the post/linux/gather/gnome_keyring_dump Metasploit module can also be used to perform this task.
