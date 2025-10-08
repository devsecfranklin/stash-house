all: slides

clean:
	@find . -name '*.pyc' | xargs rm -rf
	@find . -name '__pycache__' | xargs rm -rf
	@for trash in *.aux *.bbl *.blg *.lof *.log *.lot *.out *.pdf *.synctex.gz *.toc ; do \
		if [ -f "$$trash" ]; then \
			rm -rf $$trash ; \
		fi ; \
	done

slides: 
	cat begin.tex > stash.tex
	for i in $(shell ls -d -- section*/) ; do \
		for j in $$(ls -f $$i*.tex) ; do \
			cat $$j >> stash.tex ; \
		done; \
	done
	cat end.tex >> stash.tex
	latexmk -pdf -file-line-error -interaction=nonstopmode -synctex=1 -shell-escape stash
# bibtex stash
	latexmk -pdf -file-line-error -interaction=nonstopmode -synctex=1 -shell-escape stash
