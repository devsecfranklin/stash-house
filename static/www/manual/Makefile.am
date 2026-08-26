SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

# LaTeX lab manual build targets
# Uses xelatex for UTF-8/cjk support

.PHONY: all clean pdf

all: pdf

pdf: lab-manual.pdf

lab-manual.pdf: lab-manual.tex
	xelatex --shell-escape -interaction=nonstopmode lab-manual.tex
	xelatex --shell-escape -interaction=nonstopmode lab-manual.tex

clean:
	rm -f *.aux *.log *.toc *.out *.lot *.lof *.fls *.fdb_latexmk *.synctex.gz
	rm -f *.pdf
	find sections/ -name '*.aux' -exec rm -f {} +

dist-clean: clean
	rm -rf _build
