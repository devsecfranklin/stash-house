# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

ACLOCAL_AMFLAGS = -I config/m4 -I aclocal/latex-m4
CLOCAL_AMFLAGS = -I config/m4 -I aclocal/latex-m4

#include target.mk

check: DOCKER-exists
PYTHON-exists: ; @which python > /dev/null
mytarget: check
.PHONY: check PYTHON-exists

clean:
if HAS_DOCKER
        docker system prune -f
endif
        rm -rf *.egg-info
        find . -name '*.pyc' | xargs rm -rf
        find . -name '__pycache__' | xargs rm -rf
        for trash in AUTHORS ChangeLog aclocal.m4 autom4te.cache config config.log \
                config.status configure configure~ config.h.in~ libtool Makefile Makefile.in NEWS \
	;do \
                if [ -f $$trash ] || [ -d $$trash ]; then \
                        echo "Removing $$trash" ;\
                        rm -rf $$trash ;\
                fi ; \
        done

dist-clean: clean
        rm -rf _build

