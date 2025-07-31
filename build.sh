#!/bin/bash
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
sudo docker build -f ./docker_application_build --tag tarantula.admin:$version --build-arg app=admin --build-arg h=$host --build-arg n=admin$id --build-arg s=$seq --build-arg g=$grp .
((seq++))
sudo docker build -f ./docker_application_build --tag tarantula.presence:$version --build-arg app=presence --build-arg h=$host --build-arg n=presence$id --build-arg s=$seq --build-arg g=$grp . 
((seq++))
sudo docker build -f ./docker_application_build --tag tarantula.profile:$version --build-arg app=profile --build-arg h=$host --build-arg n=profile$id --build-arg s=$seq --build-arg g=$grp . 
((seq++))
sudo docker build -f ./docker_application_build --tag tarantula.inventory:$version --build-arg app=inventory --build-arg h=$host --build-arg n=inventory$id --build-arg s=$seq --build-arg g=$grp .
((seq++))
sudo docker build -f ./docker_application_build --tag tarantula.asset:$version --build-arg app=asset --build-arg h=$host --build-arg n=asset$id --build-arg s=$id --build-arg g=$grp . 
((seq++))
sudo docker build -f ./docker_application_build --tag tarantula.tournament:$version --build-arg app=tournament --build-arg h=$host --build-arg n=tournament$id --build-arg s=$seq --build-arg g=$grp .  
((seq++))
sudo docker build -f ./docker_application_build --tag tarantula.shop:$version --build-arg app=shop --build-arg h=$host --build-arg n=shop$id --build-arg s=$seq --build-arg g=$grp .  
sudo docker build -f ./docker_nginx_build --tag tarantula.nginx:$version .
sudo docker builder prune -af
rm id_ed25519
rm known_hosts
rm .gitconfig