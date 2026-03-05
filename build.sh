#!/bin/bash

Clean(){
    echo "deleting build files"
    rm id_ed25519
    rm known_hosts
    rm .gitconfig
    rm token.txt
}
Check(){
    if [[ $? -ne 0 ]]; then
        echo "build failed, try again"
        Clean
        exit $?
    fi
}

if [[ -n "$1" ]]; then
    version="$1"
else
    version=latest
fi

echo "Build params : ${version}"
cp ~/.ssh/id_ed25519 .
cp ~/.ssh/known_hosts .
cp ~/.gitconfig .
cp ~/token.txt .
apps=("admin" "presence" "inventory" "asset" "tournament" "postoffice" "mahjong")
for app in "${apps[@]}"; do
  echo "Current build target : $app"
  sudo docker build -f ./docker_application_build --tag tarantula.$app:$version --build-arg app=$app .
  Check
  ((seq++))
done

sudo docker build -f ./docker_prometheus_node_exporter_build --tag tarantula.node:$version .
Check

sudo docker build -f ./docker_prometheus_build --tag tarantula.prometheus:$version .
Check

sudo docker build -f ./docker_nginx_build --tag tarantula.nginx:$version .
Check
sudo docker builder prune -af

Clean
