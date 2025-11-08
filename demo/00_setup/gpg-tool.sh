#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 11/07/2025 initial version

LRED=$(tput setaf 1)

if [ -f "../../bin/common.sh" ]; then
	source "../../bin/common.sh"
else
	echo -e "${LRED}can not find common.sh${NC}"
	exit 1
fi


log_header "Show your existing keys:"
gpg --list-secret-keys --keyid-format=long

# https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key
gpg --full-generate-key
