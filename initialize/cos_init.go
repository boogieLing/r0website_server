package initialize

import (
	"fmt"
	"r0Website-server/config"
	"r0Website-server/global"
	"r0Website-server/utils"
)

// InitCOSClient 初始化COS客户端
func InitCOSClient(cfg *config.SystemConfig) error {
	cosClient, err := utils.NewCOSClient(cfg)
	if err != nil {
		return fmt.Errorf("初始化COS客户端失败: %v", err)
	}

	global.COSClient = cosClient
	return nil
}