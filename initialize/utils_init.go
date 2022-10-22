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
