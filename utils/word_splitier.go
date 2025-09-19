// Package utils
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 分词工具，对于text数组的遍历分词之后的free，可能存在内存问题！
 * 	FOLLOW: https://github.com/yanyiwu/gojieba
 *  弃用gojieba，FOLLOW: https://github.com/go-ego/gse
 * @File:  WordSplit
 * @Version: 1.0.0
 * @Date: 2022/7/5 11:53
 */
package utils

import (
	"github.com/go-ego/gse"
	"strings"
)

var (
	WordSplitSeg gse.Segmenter
)

func WordSplitForSearching(text string) string {
	words := WordSplitSeg.CutSearch(text)
	return strings.Join(words, " ")
}
