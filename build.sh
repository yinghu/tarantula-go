#!/bin/bash
version=$1
host=$2
id=$3
docker build -f ./docker_application_build --tag tarantula.admin:$version --build-arg app=admin --build-arg h=$host --build-arg n=a$id --build-arg i=$id .
docker build -f ./docker_application_build --tag tarantula.presence:$version --build-arg app=presence --build-arg h=$host --build-arg n=p$id --build-arg i=$id . 
docker build -f ./docker_application_build --tag tarantula.profile:$version --build-arg app=profile --build-arg h=$host --build-arg n=f$id --build-arg i=$id . 
docker build -f ./docker_application_build --tag tarantula.inventory:$version --build-arg app=inventory --build-arg h=$host --build-arg n=v$id --build-arg i=$id .
docker build -f ./docker_application_build --tag tarantula.asset:$version --build-arg app=asset --build-arg h=$host --build-arg n=s$id --build-arg i=$id . 
docker build -f ./docker_application_build --tag tarantula.tournament:$version --build-arg app=tournament --build-arg h=$host --build-arg n=t$id --build-arg i=$id .   