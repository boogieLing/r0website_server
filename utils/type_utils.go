package utils

import (
	"fmt"
	"r0Website-server/utils/exception"
	"reflect"
	"strconv"
	"strings"
)

// ParseBool 字符串转bool
func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	}
	return false, exception.NewSysError("转换失败")
}

func TypeChange(source string, splitStr string) []int {
	tempArr := strings.Split(source, splitStr)
	result := make([]int, 0)
	for _, item := range tempArr {
		newItem, _ := strconv.Atoi(item)
		result = append(result, newItem)
	}
	return result
}

// StructCopy 结构体复制
// source 当前有值的结构体
// target 接受值的结构体
// fields 需要的设置的属性
func StructCopy(source interface{}, target interface{}, fields ...string) (err error) {
	sourceKey := reflect.TypeOf(source)
	sourceVal := reflect.ValueOf(source)

	targetKey := reflect.TypeOf(target)
	targetVal := reflect.ValueOf(target)

	if targetKey.Kind() != reflect.Ptr {
		err = fmt.Errorf("被覆盖的数据必须是一个结构体指针")
		return
	}

	targetVal = reflect.ValueOf(targetVal.Interface())

	// 存放字段
	fieldItems := make([]string, 0)

	if len(fields) > 0 {
		fieldItems = fields
	} else {
		for i := 0; i < sourceVal.NumField(); i++ {
			fieldItems = append(fieldItems, sourceKey.Field(i).Name)
		}
	}

	for i := 0; i < len(fieldItems); i++ {
		field := targetVal.Elem().FieldByName(fieldItems[i])
		value := sourceVal.FieldByName(fieldItems[i])
		if field.IsValid() && field.Kind() == value.Kind() {
			field.Set(value)
		}
	}
	return
}
