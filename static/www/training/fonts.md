# Slides

* If you want the full power of TeX Live Manager, we [recommend installing TeX Live from upstream](https://tug.org/texlive/quickinstall.html)
* See also [Integrating vanilla TeX Live with Debian](https://tug.org/texlive/debian.html)

"Technical guidance for this project's LaTeX and Git configurations was provided by an AI assistant, powered by Google's Gemini model."

## Stuff

* [tex live](https://www.tug.org/texlive/acquire-iso.html)
* TexLive [quick install for UNIX](https://tug.org/texlive/quickinstall.html)

```sh
sudo apt-get install --reinstall texlive-base -y
sudo texhash
mkdir ~/texmf
tlmgr init-usertree
#tlmgr update --self # update
# Cross release updates are only supported with
# update-tlmgr-latest.sh --update
wget https://mirror.ctan.org/systems/texlive/tlnet/update-tlmgr-latest.sh && chmod 755 update-tlmgr-latest.sh
sh update-tlmgr-latest.sh -- --upgrade
tlmgr install montserrat
tlmgr install collection-fontsrecommended
tlmgr search --file montserrat.sty
```
