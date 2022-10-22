#!/bin/bash
unset GOPATH # GOPATH与go.mod不需要共存
go env -w GOPROXY=https://goproxy.cn,direct  #设置变量
go env #检查变量设置生效