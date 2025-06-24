# tarantula migration go version

# To build docker image 
  docker build --file[-f] cmd/web/{application}/docker_build -t[--tag] tarantula.{application}:v1 --build-arg cfg={nhost}#{nname}#{nid} .

  docker run -d --name presence --network tarantula tarantula.presence:v1.0

  docker build --file[-f] .\docker_nginx_build -t[--tag] tarantula.nginx .

  build.bat|build.sh [version] [host] [nodeid]
