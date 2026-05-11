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
apps=("admin" "presence" "inventory" "asset" "postoffice")
for app in "${apps[@]}"; do
  echo "Current build target : $app"
  podman build --no-cache -f ./docker_application_build --build-arg app=$app --tag tarantula.$app:$version .
  Check
  ((seq++))
done

podman build --no-cache  -f ./docker_nginx_build --tag tarantula.nginx:$version .
Check
podman image prune --filter "label=stage=builder"

Clean
