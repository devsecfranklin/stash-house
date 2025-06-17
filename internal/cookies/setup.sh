#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT

go mod init bitsmasher.net/cookies
go mod tidy
go mod verify
