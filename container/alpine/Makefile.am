all:
# docker image build -t stash/alpine .
        docker compose build stash
        @echo "Size of image: $(shell docker image ls| grep ^mmmbot | awk -F '  +' '{print $$5}')"
        @echo "Tagging image:  $(shell docker image ls| grep ^mmmbot | awk -F '  +' '{print $$3}')"


clean:
        @for trash in Makefile Makefile.am; do \
        if [ -f $$trash ] || [ -d $$trash ]; then \
                echo "Removing $$trash" ;\
                rm -rf $$trash ;\
        fi ; \
        done
        @if [ -d "/nix" ]; then nix-collect-garbage -d; fi
        rm -rf _build *.egg-info
        @find . -name '*.pyc' | xargs rm -rf
        @find . -name '__pycache__' | xargs rm -rf
	docker system prune -f