#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# v0.1 10/24/2025 made script for a talk demo

# SCRIPT_DIR="${0%/*}"
if [ -f "../../bin/common.sh" ]; then
	source "../../bin/common.sh"
else
	echo "can not find common.sh"
	exit 1
fi

DIRS=(${HOME}/.aws "${HOME}/.config/gcloud") # create these directories
for i in ${DIRS[@]}; do
	if [ ! -d "${i}" ]; then
		mkdir -p ${i} && log_success "Created dir: ${i}"
	else
		log_info "The directory exists: ${j}"
	fi
done

cp demo/01_locate/.files/config demo/01_locate/.files/credentials "${HOME}/.aws"
cp demo/01_locate/.files/gcp-gcs-pso-11d22411ba83.json "${HOME}/.config/gcloud"
