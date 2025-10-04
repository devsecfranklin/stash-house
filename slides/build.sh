#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
#
# ChangeLog:
#
# why did I write this, because I forgot I put it in the Makefile already.
# THat's why.

source "../bin/common.sh"

cat begin.tex > stash.tex


FOLDERS=($(ls -d */))

for i in "${FOLDERS[@]}"; do
  FILES=($(ls ${i}slide* ))
  for j in "${FILES[@]}";
  do
    cat ${j} >> stash.tex
  done
  unset FILES
done

cat end.tex >> stash.tex
