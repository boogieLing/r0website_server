//go:build windows
// +build windows

package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func initServer(port string, router *gin.Engine) server {
	// release版本
	gin.SetMode(gin.ReleaseMode)
	return &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
