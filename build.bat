@echo off
SET version=%1
IF "%version%" == "" (
    SET version=latest
)
SET host=%2
IF "%host%" == "" (
    SET host=localhost
)
SET id=%3
IF "%id%" == "" (
    SET id=1
)
SET seq=%4
IF "%seq%" == "" (
    SET /A seq=1+0
)
SET grp=%5
IF "%grp%" == "" (
    SET grp=dev
)
@echo "Build params : %version% %host% %id% %seq% %grp%"
xcopy "%HOMEDRIVE%/%HOMEPATH%"\.ssh\id_ed25519 . /s
xcopy "%HOMEDRIVE%/%HOMEPATH%"\.ssh\known_hosts . /s
xcopy "%HOMEDRIVE%/%HOMEPATH%"\.gitconfig . /s
docker build -f .\docker_application_build --tag tarantula.admin:%version% --build-arg app=admin --build-arg h=%host% --build-arg n=admin%id% --build-arg s=%seq% --build-arg g=%grp% --progress=plain .
IF %ERRORLEVEL% NEQ 0 ( 
   call :Clean
   exit 
)
SET /A seq=%seq%+1
docker build -f .\docker_application_build --tag tarantula.presence:%version% --build-arg app=presence --build-arg h=%host% --build-arg n=presence%id% --build-arg s=%seq% --build-arg g=%grp% . 
IF %ERRORLEVEL% NEQ 0 ( 
   @echo "build failed, try again"
   call :Clean
   exit 
)
SET /A seq=%seq%+1
docker build -f .\docker_application_build --tag tarantula.inventory:%version% --build-arg app=inventory --build-arg h=%host% --build-arg n=inventory%id% --build-arg s=%seq% --build-arg g=%grp% . 
IF %ERRORLEVEL% NEQ 0 (
   @echo "build failed, try again"  
   call :Clean
   exit 
)
SET /A seq=%seq%+1
docker build -f .\docker_application_build --tag tarantula.asset:%version% --build-arg app=asset --build-arg h=%host% --build-arg n=asset%id% --build-arg s=%seq% --build-arg g=%grp% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
   call :Clean
   exit 
)
SET /A seq=%seq%+1
docker build -f .\docker_application_build --tag tarantula.tournament:%version% --build-arg app=tournament --build-arg h=%host% --build-arg n=tournament%id% --build-arg s=%seq% --build-arg g=%grp% . 
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
   call :Clean
   exit 
)
SET /A seq=%seq%+1
docker build -f .\docker_application_build --tag tarantula.postoffice:%version% --build-arg app=postoffice --build-arg h=%host% --build-arg n=postoffice%id% --build-arg s=%seq% --build-arg g=%grp% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
   call :Clean
   exit 
)
SET /A seq=%seq%+1
docker build -f .\docker_application_build --tag tarantula.mahjong:%version% --build-arg app=mahjong --build-arg h=%host% --build-arg n=mahjong%id% --build-arg s=%seq% --build-arg g=%grp% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
   call :Clean
   exit 
)
docker build -f .\docker_nginx_build --tag tarantula.nginx:%version% .
IF %ERRORLEVEL% NEQ 0 ( 
    @echo "build failed, try again"
   call :Clean
   exit 
)
docker builder prune -af
call :Clean

:Clean
@echo "deleting build files"
del id_ed25519
del known_hosts
del .gitconfig
exit /B 0