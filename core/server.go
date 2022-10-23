// Package core
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 服务器核心
 * @File:  server
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:35
 */
package core

import (
	"r0Website-server/global"
	"r0Website-server/router"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	// 初始化路由
	Router := router.Routers()
	// 初始化服务
	s := initServer(global.Config.System.Port, Router)
	// 启动服务器
	s.ListenAndServe().Error()
}
