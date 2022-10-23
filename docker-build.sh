#!/bin/bash


# wget https://golang.google.cn/dl/go1.17.11.linux-amd64.tar.gz -P ./
docker container stop r0website_server; docker container rm r0website_server ; docker image rm r0website_server;
# docker network rm r0-work-net;

# 以网桥模式新建一个docker网络
docker network create -d bridge --subnet 172.97.0.0/16 --gateway 172.97.0.1 r0-work-net

docker build -f ./Dockerfile --target server_apply -t r0website_server . &&
docker run --name r0website_server -d -p 8202:8202 --network r0-work-net r0website_server:latest

docker container ls -a

# chmod +x ./docker-build.sh && ./docker-build.sh