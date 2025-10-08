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
                config.status configure configure~ config.h.in~ libtool Makefile Makefile.in NEWS ;do \
                if [ -f $$trash ] || [ -d $$trash ]; then \
                        echo "Removing $$trash" ;\
                        rm -rf $$trash ;\
                fi ; \
        done

dist-clean: clean
        rm -rf _build

function check_python_version() {
  echo -e "\n${LPURP}# --- Check Python Version -------------------------------------------\n${NC}" | tee -a "${RAW_OUTPUT}"
  if command -v python &>/dev/null; then
    PYTHON_VERSION=$(python -c 'import sys; print(sys.version_info.major)')
    if [[ "$PYTHON_VERSION" -eq 3 ]]; then
      echo -e "${LBLUE}The 'python' command points to Python 3.${NC}"
      # Use 'python'
      PYTHON_CMD="python"
    elif [[ "$PYTHON_VERSION" -eq 2 ]]; then
      echo -e "${LBLUE}The 'python' command points to Python 2.${NC}"

      # Decide what to do: try python3, or exit
      if command -v python3 &>/dev/null; then
        echo -e "Using 'python3' instead."
        PYTHON_CMD="python3"
      else
        echo -e "${LRED}Error: Python 3 not found. Exiting.${NC}"
        exit 1
      fi
    else
      echo -e "${LRED}The 'python' command points to an unknown Python version ($PYTHON_VERSION).${NC}"
      # Decide what to do
      if command -v python3 &>/dev/null; then
        echo "Attempting to use 'python3' instead.${NC}"
        PYTHON_CMD="python3"
      else
        echo -e "${LRED}Error: Python 3 not found. Exiting.${NC}"
        exit 1
      fi
    fi
  elif command -v python3 &>/dev/null; then
    echo "${LGREEN}Using python3 as your python version.${NC}"
    PYTHON_CMD="python3"
  else
    echo -e "${LRED}Error: Neither 'python' nor 'python3' found. Please install Python 3. Exiting.${NC}"
    exit 1
  fi

  echo -e "${LBLUE}Using Python command: ${PYTHON_CMD}${NC}"
}
