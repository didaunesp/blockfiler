#!/bin/bash

main(){
  rm -rf ./scripts/users/*
  clearDocker
}

clearDocker(){
  docker kill $(docker ps -aq)
  docker rm $(docker ps -aq)
  docker network prune -f
  docker volume prune -f
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*chaincode*/) {print $3}')
    if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
      echo "---- No images available for deletion ----"
    else
      docker rmi -f $DOCKER_IMAGE_IDS
    fi
}

main