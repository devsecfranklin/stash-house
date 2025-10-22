#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# ChangeLog:

if [ ! -d "aclocal" ]; then
  echo "create aclocal dir"
  mkdir -p aclocal
fi

if [ ! -d "config/m4" ]; then
  echo "create congif/m4 dir"
  mkdir -p config/m4
fi

if [ ! -f "Makefile.in" ] && [ -f "./config.status" ]; then
  rm config.status # if Makefile.in is missing, then erase stale config.status
fi

aclocal -Iaclocal/latex-m4 || exit 1
autoreconf -i
automake -a -c --add-missing || exit 1
