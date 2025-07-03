#!/bin/bash
version=$1
host=$2
id=$3
seq=$4
if [[ -n "$5" ]]; then
    grp="$5"
else
    grp="dev"
fi
docker build -f ./docker_application_build --tag tarantula.admin:$version --build-arg app=admin --build-arg h=$host --build-arg n=admin$id --build-arg i=$seq --build-arg g=$grp .
((seq++))
docker build -f ./docker_application_build --tag tarantula.presence:$version --build-arg app=presence --build-arg h=$host --build-arg n=presence$id --build-arg i=$seq --build-arg g=$grp . 
((seq++))
docker build -f ./docker_application_build --tag tarantula.profile:$version --build-arg app=profile --build-arg h=$host --build-arg n=profile$id --build-arg i=$seq --build-arg g=$grp . 
((seq++))
docker build -f ./docker_application_build --tag tarantula.inventory:$version --build-arg app=inventory --build-arg h=$host --build-arg n=inventory$id --build-arg i=$seq --build-arg g=$grp .
((seq++))
docker build -f ./docker_application_build --tag tarantula.asset:$version --build-arg app=asset --build-arg h=$host --build-arg n=asset$id --build-arg i=$id --build-arg g=$grp . 
((seq++))
docker build -f ./docker_application_build --tag tarantula.tournament:$version --build-arg app=tournament --build-arg h=$host --build-arg n=tournament$id --build-arg i=$seq --build-arg g=$grp .   
docker build -f ./docker_nginx_build --tag tarantula.nginx:$version .