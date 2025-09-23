// Package initialize
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 初始化工具相关配置
 * @File:  utils
 * @Version: 1.0.0
 * @Date: 2022/7/5 13:41
 */
package initialize

import (
	"r0Website-server/global"
	"r0Website-server/utils"
)

func InitUtils() {
	global.Logger.Infoln("InitUtils")
	initWordSplitSeg()
	// COS初始化移到R0Ioc初始化中
}
func initWordSplitSeg() {
	var err error
	err = utils.WordSplitSeg.LoadDict("zh_s")
	if err != nil {
		global.Logger.Error(err)
	}
	err = utils.WordSplitSeg.LoadDictEmbed("zh_s")
	if err != nil {
		global.Logger.Error(err)
	}
}

func initCOSClient() {
	// 检查配置是否已初始化
	if global.Config == nil {
		global.Logger.Error("全局配置未初始化，跳过COS客户端初始化")
		return
	}

	// 初始化COS客户端
	if err := InitCOSClient(global.Config); err != nil {
		global.Logger.Errorf("初始化COS客户端失败: %v", err)
		return
	}
	global.Logger.Infoln("COS客户端初始化成功")
}
