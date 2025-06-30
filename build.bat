@echo off
SET version=%1
SET host=%2
SET id=%3
SET seq=%4
docker build -f .\docker_application_build --tag tarantula.admin:%version% --build-arg app=admin --build-arg h=%host% --build-arg n=admin%id% --build-arg i=%seq% .
SET /A seq=seq+1
docker build -f .\docker_application_build --tag tarantula.presence:%version% --build-arg app=presence --build-arg h=%host% --build-arg n=presence%id% --build-arg i=%seq% . 
SET /A seq=seq+1
docker build -f .\docker_application_build --tag tarantula.profile:%version% --build-arg app=profile --build-arg h=%host% --build-arg n=profile%id% --build-arg i=%seq% .  
SET /A seq=seq+1
docker build -f .\docker_application_build --tag tarantula.inventory:%version% --build-arg app=inventory --build-arg h=%host% --build-arg n=inventory%id% --build-arg i=%seq% . 
SET /A seq=seq+1
docker build -f .\docker_application_build --tag tarantula.asset:%version% --build-arg app=asset --build-arg h=%host% --build-arg n=asset%id% --build-arg i=%seq% .
SET /A seq=seq+1
docker build -f .\docker_application_build --tag tarantula.tournament:%version% --build-arg app=tournament --build-arg h=%host% --build-arg n=tournament%id% --build-arg i=%seq% . 
docker build -f .\docker_nginx_build --tag tarantula.nginx:%version% .