// Package dao
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description:
 * @File:  user_dao
 * @Version: 1.0.0
 * @Date: 2022/7/31 18:10
 */
package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"r0Website-server/global"
	"r0Website-server/models/po"
)

type UserDao struct {
	*BasicDaoMongo `R0Ioc:"true"`
}

func (*UserDao) CollectionName() string {
	return "user"
}
func (ud *UserDao) Collection() *mongo.Collection {
	return ud.Mdb.Collection(ud.CollectionName())
}

// FindObjByEmail 通过email查找用户
func (ud *UserDao) FindObjByEmail(email string) (po.User, error) {
	filter := bson.M{"email": email}
	var ans po.User
	err := ud.Collection().FindOne(context.Background(), &filter).Decode(&ans)
	if err != nil {
		global.Logger.Error(err)
	}
	return ans, err
}

// CreateUser 新增一个用户
func (ud *UserDao) CreateUser(input po.User) error {
	_, err := ud.Collection().InsertOne(context.TODO(), input)
	if err != nil {
		global.Logger.Error(err)
	}
	return err
}

// EmailCount 统计拥有此Email的用户数量
func (ud *UserDao) EmailCount(email string) int64 {
	filter := bson.M{"email": email}
	count, err := ud.Collection().CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Logger.Error(err)
	}
	return count
}
