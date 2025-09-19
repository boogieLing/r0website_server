/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 提供网站的后端网络服务，服务器启动入口
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:16
 */
package main

import (
	"fmt"
	"os"
	"r0Website-server/core"
	"r0Website-server/global"
	"r0Website-server/initialize"
	"r0Website-server/r0Ioc"
	"strconv"
	"strings"
)

// using env:   export GIN_MODE=release
func main() {
	// 解析参数
	if len(os.Args) > 1 {
		for idx, arg := range os.Args {
			fmt.Println("参数"+strconv.Itoa(idx)+" : ", arg)
		}
		arg := strings.Split(os.Args[1], "=")
		if len(arg) >= 2 && arg[0] == "--config" {
			// 指定yaml文件的路径
			global.Config = initialize.InitProdConfig(arg[1])
		}
	} else {
		global.Config = initialize.InitDevConfig()
	}
	initialize.InitLogger()
	initialize.InitUtils()
	r0Ioc.InitR0Ioc(global.YAMLPATH)
	defer r0Ioc.ExitR0Ioc()
	core.RunWindowsServer()
}
