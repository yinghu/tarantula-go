#!/bin/bash

Clean(){
    echo "deleting build files"
    rm id_ed25519
    rm known_hosts
    rm .gitconfig
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
if [[ -n "$2" ]]; then
    host="$2"
else
    host="localhost"
fi
if [[ -n "$3" ]]; then
    id="$3"
else
    id=1
fi
if [[ -n "$4" ]]; then
    seq="$4"
else
    seq=1
fi
if [[ -n "$5" ]]; then
    grp="$5"
else
    grp="dev"
fi
echo "Build params : ${version} ${host} ${id} ${seq} ${grp}"
cp ~/.ssh/id_ed25519 .
cp ~/.ssh/known_hosts .
cp ~/.gitconfig .
apps=("admin" "presence" "inventory" "asset" "tournament" "postoffice" "mahjong")
for app in "${apps[@]}"; do
  echo "Current build target : $app"
  sudo docker build -f ./docker_application_build --tag tarantula.$app:$version --build-arg app=$app --build-arg h=$host --build-arg n=$app$id --build-arg s=$seq --build-arg g=$grp .
  Check
  ((seq++))
done

sudo docker build -f ./docker_nginx_build --tag tarantula.nginx:$version .
Check
sudo docker builder prune -af

Clean
