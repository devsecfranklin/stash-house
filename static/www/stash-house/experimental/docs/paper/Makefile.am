# SPDX-FileCopyrightText: 2026 franklin <franklin.diaz@ucdenver.edu>
#
# SPDX-License-Identifier: MIT

clean:
if HAS_DOCKER
        docker system prune -f
endif
	latexmk -C
	rm -f *.bbl *.bcf *.run.xml *.nav *.snm *.vrb

all:
	pdflatex stash.tex
	biber stash
	pdflatex stash.tex
	pdflatex stash.tex
	latexmk -pdf -bibtex -f stash.tex
