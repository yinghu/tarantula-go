@echo off
SET version=%1
IF "%version%" == "" (
    SET version=latest
)
@echo "Build params : %version%"
xcopy "%HOMEDRIVE%/%HOMEPATH%"\.ssh\id_ed25519 . /s
xcopy "%HOMEDRIVE%/%HOMEPATH%"\.ssh\known_hosts . /s
xcopy "%HOMEDRIVE%/%HOMEPATH%"\.gitconfig . /s
xcopy "%HOMEDRIVE%/%HOMEPATH%"\token.txt . /s
docker build -f .\docker_application_build --tag tarantula.admin:%version% --build-arg app=admin .
IF %ERRORLEVEL% NEQ 0 ( 
   @echo "build failed, try again"
   goto Clean 
)

docker build -f .\docker_application_build --tag tarantula.presence:%version% --build-arg app=presence . 
IF %ERRORLEVEL% NEQ 0 ( 
   @echo "build failed, try again"
   goto Clean
)

docker build -f .\docker_application_build --tag tarantula.inventory:%version% --build-arg app=inventory . 
IF %ERRORLEVEL% NEQ 0 (
   @echo "build failed, try again"  
   goto Clean
)

docker build -f .\docker_application_build --tag tarantula.asset:%version% --build-arg app=asset .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
    goto Clean
)
docker build -f .\docker_application_build --tag tarantula.postoffice:%version% --build-arg app=postoffice .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
    goto Clean
)
docker build -f .\docker_application_build --tag tarantula.mahjong:%version% --build-arg app=mahjong .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
    goto Clean
)
docker build -f .\docker_prometheus_node_exporter_build --tag tarantula.node:%version% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
    goto Clean
)
docker build -f .\docker_prometheus_build --tag tarantula.prometheus:%version% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
    goto Clean
)
docker build -f .\docker_nginx_build --tag tarantula.nginx:%version% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
    goto Clean
)

:Clean
docker builder prune -af
@echo "deleting build files"
del id_ed25519
del known_hosts
del .gitconfig
del token.txt
