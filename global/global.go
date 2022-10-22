// Package global
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 全局使用的组件或者变量
 * @File:  global
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:20
 */
package global

import (
	"github.com/sirupsen/logrus"
	"r0Website-server/config"
)

var (
	Config *config.SystemConfig
	Logger *logrus.Logger
)
