// Package initialize
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 日志工具初始化
 * @File:  logger
 * @Version: 1.0.0
 * @Date: 2022/7/5 13:49
 */
package initialize

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"r0Website-server/global"
	"r0Website-server/utils"
	"strings"
	"time"
)

func InitLogger() {
	ioWriter := utils.YamlLogFile()
	global.Logger = logrus.New()
	global.Logger.SetOutput(ioWriter)
	global.Logger.SetLevel(logrus.DebugLevel)
	//设置日志格式
	global.Logger.SetFormatter(new(GlobalLogFormatter))
}

type GlobalLogFormatter struct{}

func (box *GlobalLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("<GLOBAL:%s> : [%s] - [%s]\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}
