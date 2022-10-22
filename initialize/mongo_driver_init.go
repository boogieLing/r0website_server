// Package initialize
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: mongo操作驱动
 * @File:  mongo_driver
 * @Version: 1.0.0
 * @Date: 2022/7/3 23:52
 */
package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"r0Website-server/config"
	"r0Website-server/dao"
	"r0Website-server/global"
	"time"
)

// InitDB 获取一个mongo的db操作实例
// 一定要记得defer关闭client
func InitDB() (context.Context, *mongo.Client) {
	// FOLLOW:
	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/#std-label-golang-authentication-mechanisms
	credential := options.Credential{
		AuthSource: global.Config.Mongo.DB,
		Username:   global.Config.Mongo.SuperAdminName,
		Password:   global.Config.Mongo.SuperAdminPswd,
	}
	dsn := fmt.Sprintf("mongodb://%s:%s/?authSource=admin",
		global.Config.Mongo.Address,
		global.Config.Mongo.Part,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn).SetAuth(credential))
	if err != nil {
		global.Logger.Error(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		global.Logger.Error(err)
	}
	// db := client.Database(global.Config.Mongo.DB)
	global.Logger.Infof("=== Connected MongoDB:%s ===", global.Config.Mongo.DB)
	return ctx, client
}

// MongoConstructor Mongo构造
func MongoConstructor(cfg *config.Mongo) *dao.BasicDaoMongo {
	credential := options.Credential{
		AuthSource: cfg.DB,
		Username:   cfg.SuperAdminName,
		Password:   cfg.SuperAdminPswd,
	}
	dsn := fmt.Sprintf("mongodb://%s:%s/?authSource=admin",
		cfg.Address,
		cfg.Part,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn).SetAuth(credential))
	if err != nil {
		global.Logger.Error(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		global.Logger.Error(err)
	}
	// db := client.Database(cfg.DB)
	global.Logger.Infof("=== Connected MongoDB:%s ===", cfg.DB)
	return &dao.BasicDaoMongo{Mc: client, Mdb: client.Database(cfg.DB)}
}
