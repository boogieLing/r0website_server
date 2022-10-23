#!/bin/bash
mkdir logs
touch logs/r0Website_server.log
unset GOPATH # GOPATH与go.mod不需要共存
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go build -o ./r0Website-server ./main/main.go
./r0Website-server

# chmod +x ./run-server.sh && ./run-server.sh