// Package utils
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 将字符串转为24位的hex串,
 * @File:  string2hexstring24
 * @Version: 1.0.0
 * @Date: 2022/7/6 01:46
 */
package utils

import (
	"encoding/hex"
	"strings"
)

func String2HexString24(input string) string {
	if len(input) == 24 {
		// 如果刚好是24，姑且认为输入的就是mongo的id，不做处理
		return input
	}
	var output string
	ans := hex.EncodeToString([]byte(input))
	output = ans
	if len(output) < 24 {
		prefix := strings.Repeat("0", 24-len(output))
		output = prefix + output
	}
	if len(output) > 24 {
		output = output[:24]
	}
	return output
}
