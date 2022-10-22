#!/bin/bash
unset GOPATH # GOPATH与go.mod不需要共存
go build -o ./r0Website-server
./r0Website-server