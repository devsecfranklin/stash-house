#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 02/25/2025 initial version

DEB_PKG=(curl git gnupg2 keyringer make pass podman-compose)
#SCRIPT_DIR="${0%/*}" # echo "$SCRIPT_DIR"

if [ -f "./common.sh" ]; then
	source "./common.sh"
else
	echo -e "${LRED}can not find common.sh, run this script from the bin/ directory.${NC}"
	exit 1
fi

function main() {
	log_header "Installing Development Environment for the Stash House"
	check_python_version
	install_debian
}

main "$@"
