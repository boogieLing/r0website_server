#!/bin/bash

# wget https://golang.google.cn/dl/go1.17.11.linux-amd64.tar.gz -p ./
docker container stop r0website_server; docker container rm r0website_server ; docker image rm r0website_server;

docker build -f ./Dockerfile --target server_apply -t r0website_server . &&
docker run --name r0website_server -d -p 8202:8202 r0website_server:latest

# chmod +x ./docker-build.sh && ./docker-build.sh