# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

ACLOCAL_AMFLAGS = -I config/m4 -I aclocal/latex-m4
CLOCAL_AMFLAGS = -I config/m4 -I aclocal/latex-m4

slides: 
	latexmk -pdf -file-line-error -interaction=nonstopmode -synctex=1 -shell-escape stash
# bibtex stash
	latexmk -pdf -file-line-error -interaction=nonstopmode -synctex=1 -shell-escape stash

