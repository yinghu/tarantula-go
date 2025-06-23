# tarantula migration go version

# To build docker image 
  docker build --file cmd/web/{application}/docker_build -t tarantula.{application}:v1 .

  docker run -d --name presence --network tarantula tarantula.presence:v1.0

  docker build --file .\docker_nginx_build -t tarantula.nginx .