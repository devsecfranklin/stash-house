#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 02/25/2025 initial version

DEB_PKG=(git gnupg keyringer pass)
LOCAL_STASH="${HOME}/.stash-house"
#SCRIPT_DIR="${0%/*}" # echo "$SCRIPT_DIR"
LRED=$(tput setaf 1)

if [ -f "./common.sh" ]; then
	source "./common.sh"
else
	echo -e "${LRED}can not find common.sh, run this script from the bin/ directory.${NC}"
	exit 1
fi

function install_local() {

	log_header "Install the local passwd storage folder"

	if [ ! -d "${LOCAL_STASH}" ]; then
		mkdir -p "${LOCAL_STASH}"
		log_info "${LOCAL_STASH} directory was created."
	else
		log_info "${LOCAL_STASH} directory already exists."
	fi

	if [ ! -f "${LOCAL_STASH}/.gitattributes" ]; then
		echo "*.gpg diff=gpg" >"${LOCAL_STASH}/.gitattributes"
		log_info "created file: ${LOCAL_STASH}/.gitattributes"
	fi

	check_installed pass
	check_installed keyringer
}

declare -a KEYSERVERS=(
	"hkp://keyserver.ubuntu.com:80"
	"keyserver.ubuntu.com"
	"ha.pool.sks-keyservers.net"
	"hkp://ha.pool.sks-keyservers.net:80"
	"p80.pool.sks-keyservers.net"
	"hkp://p80.pool.sks-keyservers.net:80"
	"pgp.mit.edu"
	"hkp://pgp.mit.edu:80"
)

GPG_DIR="${HOME}/.gnupg"
function gpg_setup() {
	log_header "Configure GnuPG"

	check_installed gnupg
	if [ ! -d "${GPG_DIR}" ]; then mkdir -p "${GPG_DIR}" && log_success "Created dir: ${GPG_DIR}"; fi

	chown -R "$(whoami):engr" "${GPG_DIR}"
	chmod 750 "${GPG_DIR}"

	log_info "Importing my public key to the work folder"
	gpg --homedir "${GPG_DIR}" --import "${GPG_DIR}/franklin.gpg"
}

function update_keys() {
	SIG=$1
	MY_FILE=$2

	log_header "update_keys from ${SIG}"

	# gpg --homedir "${PREFIX}/gnupg" --import "${PREFIX}/m4-latest.tar.xz.sig" # import the key into your local keyring
	# gpg --homedir "${BUILgpg --homedir "${PREFIX}/gnupg" D_DIR}/gnupg" --fingerprint # view the fingerprints in your local keyring
	MY_KEY="$("${SIG}" 2>&1 | grep RSA | awk '{print $5}')"

	# check if the MY_KEY string is blank
	if [ -n "${MY_KEY}" ]; then
		echo -e "${LGREEN}Found key: ${MY_KEY}${NC}"
	else
		echo -e "${LRED}No key found in ${SIG}${NC}"
		return
	fi

	for server in "${KEYSERVERS[@]}"; do
		# check if MY_KEY already exists locally
		#LOCAL_KEY=$(gpg --homedir ${PREFIX}/gnupg --export-options export-minimal --armor --export ${MY_KEY} 2>&1)
		LOCAL_KEY="$(
			gpg --homedir ${PREFIX}/gnupg --export-options export-minimal --armor --export "${MY_KEY}" 2>&1
		)"
		if [ -n "${LOCAL_KEY}" ]; then
			echo -e "${CYAN}Local copy of key not found, getting${NC}"
			gpg --homedir "${PREFIX}/gnupg" --keyserver "${server}" --recv-keys "${MY_KEY}"
			if [ $? -eq 0 ]; then
				echo -e "${LGREEN}Success importing key!"
				return
			else
				echo -e "${LRED}Not finding this key: ${MY_KEY} at server: ${server}${NC}"
			fi
		fi

		VERIFIED=$(gpg --homedir "${PREFIX}/.gnupg" --verify "${SIG}" "${MY_FILE}" 2>&1 | grep 'Good signature')
		if [[ "$VERIFIED" ]]; then
			echo -e "${LGREEN}gpg keys verified. Installing...${NC}"
		else
			echo -e "${LRED}gpg key cannot be verifed. Aborting installation of ${MY_FILE}${NC}"
			exit 1
		fi
	done
}

function main() {
	log_header "Installing Stash House"
	check_if_root
	check_container
	install_local
	# install_debian
}

main "$@"
