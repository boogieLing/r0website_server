//go:build !windows
// +build !windows

package core

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"r0Website-server/utils"
	"strconv"
	"syscall"
	"time"
)

func initServer(port string, router *gin.Engine) server {
	gin.SetMode(gin.ReleaseMode)
	s := endless.NewServer(port, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 60 * time.Second
	s.MaxHeaderBytes = 1 << 20
	s.BeforeBegin = func(_ string) {
		fmt.Printf("r0website server in %d:%s\n", syscall.Getpid(), port)
		utils.SetPid(strconv.Itoa(syscall.Getpid()))
	}
	return s
}
