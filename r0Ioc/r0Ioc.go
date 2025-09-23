// Package r0Ioc
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 依赖注入容器
 * @File:  R0Ioc
 * @Version: 1.0.0
 * @Date: 2022/7/31 01:15
 */
package r0Ioc

import (
	"fmt"
	"r0Website-server/api/admin"
	"r0Website-server/api/base"
	"r0Website-server/dao"
	"r0Website-server/global"
	"r0Website-server/initialize"
	"r0Website-server/utils"
	"reflect"
)

// RegisterComponents 注册组件
func RegisterComponents(components ...interface{}) {
	for _, v := range components {
		RegisterComponentSingle(v, nil)
	}
}

// RegisterComponentSingle 注册单个组件，并允许组件有结束方法
func RegisterComponentSingle(component interface{}, exitFunc func(item *R0IocItem)) {
	componentsKey := reflect.TypeOf(component)
	componentsVal := reflect.ValueOf(component)
	if componentsKey.Kind() != reflect.Ptr {
		fmt.Printf("[R0Ioc] ERROR: %s is %T not a ptr\n", componentsKey.String(), componentsKey.Kind())
		return
	}
	componentsName := componentsKey.Elem().Name()          // 获取的是实际的类型, 比如 BasicDao
	componentsType := componentsVal.Elem().Type().String() // 获取的是导入实例的类型，比如 dao.BasicDao
	R0Ioc[componentsName] = &R0IocItem{
		Name:     componentsName,
		Type:     componentsType,
		Instance: component,
		ExitFunc: exitFunc,
	}
	if R0IocDebug {
		fmt.Printf("[R0Ioc RegisterComponents]: %s, %s\n", componentsName, componentsType)
	}

}

// Injection 注入器，将组件注入到需要的结构体
func Injection(targets ...interface{}) {
	for _, v := range targets {
		dfsInjection(reflect.TypeOf(v), reflect.ValueOf(v), make(map[reflect.Type]bool))
	}
}
func dfsInjection(refKey reflect.Type, refVal reflect.Value, visited map[reflect.Type]bool) {
	if _, ok := visited[refKey]; ok {
		return
	}
	switch refKey.Kind() {
	case reflect.Ptr:
		refKey = refKey.Elem()
		refVal = refVal.Elem()
		for i := 0; i < refKey.NumField(); i++ {
			fKey := refKey.Field(i)
			fVal := refVal.Field(i)
			fKeyName := fKey.Type.Elem().Name()
			if R0IocDebug {
				fmt.Printf("[R0Ioc dfsInjection]: %s,%s\n", fKey.Type, fVal.String())
				fmt.Printf("Cur fkeyName: %s\n", fKeyName)
			}
			if components, ok := R0Ioc[fKeyName]; ok && fKey.Tag.Get("R0Ioc") == "true" {
				// 匹配到注入目标，注入
				if R0IocDebug {
					// 正常的情况 Yes时必定是一个 <*r0Ioc.R0IocItem Value>
					fmt.Printf("Yes: ")
					fmt.Println(reflect.ValueOf(components).String())
				}
				fVal.Set(reflect.ValueOf(components.Instance))
			} else if fKey.Tag.Get("R0Ioc") == "true" {
				if R0IocDebug {
					fmt.Printf("No: ")
					fmt.Println(fVal.Type().Elem())
				}
				tmp := reflect.New(fVal.Type().Elem())
				if m := tmp.MethodByName("Construct"); m.IsValid() {
					m.Call(nil)
				}
				fVal.Set(tmp)
				visited[refKey] = true
				dfsInjection(fKey.Type, fVal, visited)
			}
		}
	default:
		fmt.Printf("[R0Ioc] ERROR: %s is %T not a ptr\n", refKey.String(), refKey.Kind())
		return
	}
}

// InitR0Route 初始化路由
func InitR0Route() []interface{} {
	var res []interface{}
	refKey := reflect.TypeOf(R0Route).Elem()
	refVal := reflect.ValueOf(R0Route).Elem()
	for i := 0; i < refKey.NumField(); i++ {
		fVal := refVal.Field(i)
		// 获取类型并创建，注意还需要一个Elem()
		tmp := reflect.New(fVal.Type().Elem())
		fVal.Set(tmp)
		// 能够成功，内存哲学
		res = append(res, tmp.Interface())
	}
	return res
}

const R0IocDebug = false

type R0IocItem struct {
	Name     string
	Type     string
	Instance interface{}
	ExitFunc func(item *R0IocItem)
}

var R0Ioc = map[string]*R0IocItem{}

var R0Route = &struct {
	AdminCategoryController   *admin.CategoryController
	AdminArticleController    *admin.ArticleController
	AdminUserController       *admin.UserController
	BaseCategoryController    *base.CategoryController
	BaseArticleController     *base.ArticleController
	BaseUserController        *base.UserController
	PicBedAlbumController     *base.PicBedAlbumController
	PicBedImageController     *base.PicBedImageController
	ImageCategoryController   *base.ImageCategoryController
	TagController             *base.TagController
}{}

// InitR0Ioc 初始化容器
func InitR0Ioc(configPath string) {
	cfg := initialize.InitProdConfig(configPath)
	basicDao := initialize.MongoConstructor(&cfg.Mongo)

	// 创建COS客户端
	var cosClient *utils.COSClient
	if err := initialize.InitCOSClient(cfg); err != nil {
		fmt.Printf("初始化COS客户端失败: %v\n", err)
		// 不中断程序，但COS功能将不可用
		cosClient = nil
	} else {
		cosClient = global.COSClient
	}

	RegisterComponents([]interface{}{
		cfg,
		cosClient,
	}...)
	RegisterComponentSingle(basicDao, func(item *R0IocItem) {
		item.Instance.(*dao.BasicDaoMongo).Disconnect()
	})
	Injection(InitR0Route()...)
	for key, value := range R0Ioc {
		fmt.Println("[R0IOC]\t", key, "\t---->\t", value)
	}
}

// ExitR0Ioc 退出容器
func ExitR0Ioc() {
	for _, value := range R0Ioc {
		if value.ExitFunc != nil {
			value.ExitFunc(value)
		}
	}
}
