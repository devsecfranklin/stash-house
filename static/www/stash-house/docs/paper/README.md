# paper

- `-pdf: generates pdf`
- `-bibtex: ensures bibliography is processed (supports biber)`
- `-f: force completion even with minor errors`

## build

```sh
sudo apt-get install texlive-pictures texlive-latex-extra
mkdir -p ~/texmf/tex/latex
tlmgr init-usertree
tlmgr install pgf
tlmgr install lipsum
tlmgr list --installed
latexmk -c
ls -lart
latexmk -pdf stash.tex
```
