// Package utils
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 将md文件内容转为html字符串
 * 	FOLLOW:https://github.com/russross/blackfriday
 * @File:  Markdown2html
 * @Version: 1.0.0
 * @Date: 2022/7/4 22:10
 */
package utils

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func Markdown2Html(md string) string {
	unsafe := blackfriday.Run(
		[]byte(md),
		blackfriday.WithExtensions(blackfriday.CommonExtensions|
			blackfriday.HardLineBreak|
			blackfriday.AutoHeadingIDs|
			blackfriday.Autolink,
		),
	)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}
