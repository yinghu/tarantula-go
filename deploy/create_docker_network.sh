#!/bin/bash
docker network create --driver bridge --subnet 10.20.0.0/16 --gateway 10.20.0.1 tarantula-net