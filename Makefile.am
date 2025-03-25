CLOCAL_AMFLAGS = -I config/m4 -I aclocal
ACLOCAL_AMFLAGS = -I config/m4 -I aclocal

clean:
	@for trash in _build aclocal aclocal.m4 autom4te.cache config config.in config.in~ config.log config.status configure configure~ libtool Makefile Makefile.in container/Makefile container/Makefile.in po stamp-h1 terraform-install.sh; do \
	if [ -f $$trash ] || [ -d $$trash ]; then \
		echo "Removing $$trash" ;\
		rm -rf $$trash ;\
	fi ; \
	done
	@if [ -d "/nix" ]; then nix-collect-garbage -d; fi
	rm -rf _build *.egg-info
	@find . -name '*.pyc' | xargs rm -rf
	@find . -name '__pycache__' | xargs rm -rf

node:
	npm install @linthtml/linthtml --save-dev
	npm install js-beautify
	@if [ ! -e ".linthtmlrc.json" ]; then ./node_modules/.bin/linthtml --init; fi
	npm ci
