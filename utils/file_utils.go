// Package utils
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 文件相关工具
 * @File:  file_utils
 * @Version: 1.0.0
 * @Date: 2022/7/5 14:15
 */
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// todayFilename 自动获取文件名
func todayFilename() string {
	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := currentPath + "/logs/" + time.Now().Format("2006-01")
	_ = os.MkdirAll(path, os.ModePerm)
	day := time.Now().Format("02")
	return path + "/" + day + ".log"
}

// NewLogFile 新建当日文件对象
func NewLogFile() *os.File {
	filename := todayFilename()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println(filename)
	}
	return file
}

// YamlLogFile yaml中配置的文件
func YamlLogFile(logPath string) *os.File {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("err:", err)
	}
	return file
}

// SetPid 进程号
func SetPid(pid string) {
	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	pidFile := currentPath + "/.pid"
	file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		n, _ := file.Seek(0, os.SEEK_END)
		_, err = file.WriteAt([]byte(pid), n)
		fmt.Printf("pid in %s\n", pidFile)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("close err:", err)
			}
		}(file)
	}
}
