// Package bo
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 各种错误状态的描述
 * @File:  error
 * @Version: 1.0.0
 * @Date: 2022/7/4 17:17
 */
package bo

import "fmt"

type UniqueError struct {
	UniqueField string
	Msg         string
	Count       int64
}

func (a *UniqueError) Error() string {
	return fmt.Sprintf("字段: %s 值已存在，但它应是唯一的!, 值[%s] count: %d", a.UniqueField, a.Msg, a.Count)
}
