#!/bin/bash

RED='\033[0;31m'
LRED='\033[1;31m'
LGREEN='\033[1;32m'
CYAN='\033[0;36m'
LPURP='\033[1;35m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function check_if_root() {
  if [[ $EUID -eq 0 ]]; then
    echo -e "${LGREEN}Superuser, good.${NC}"
  else
    echo -e "${RED}Need to be Superuser to install packages.${NC}"
    exit 1
  fi
}

function install-deps {
    apt install nginx
}

function config_proxy {
    unlink /etc/nginx/sites-enabled/default

    echo -e "${CYAN}Writing proxy config file.${NC}"
cat <<EOF > /etc/nginx/sites-available/reverse-proxy.conf
server {
    listen 8080;
    location /scustomer-proxy/ {
        proxy_pass http://httpbin.org;
    }
}
EOF
}

function main {
  check_if_root
  install_deps
  config_proxy
}

main