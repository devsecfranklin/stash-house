# SPDX-FileCopyrightText: 2026 franklin <franklin.diaz@ucdenver.edu>
#
# SPDX-License-Identifier: MIT

clean:
if HAS_DOCKER
        docker system prune -f
endif
        for trash in *.aux \
        do \
                if [ -f $$trash ] || [ -d $$trash ]; then \
                        echo "Removing $$trash" ;\
                        rm -rf $$trash ;\
                fi ; \
        done

