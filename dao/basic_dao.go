// Package dao
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description:
 * @File:  basic_dao
 * @Version: 1.0.0
 * @Date: 2022/7/30 03:05
 */
package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type BasicDaoInter interface {
	CollectionName() string
	Collection() *mongo.Collection
	Disconnect()
}
type BasicDaoMongo struct {
	Mc  *mongo.Client
	Mdb *mongo.Database
}

func (bd *BasicDaoMongo) CollectionName() string {
	return ""
}

func (bd *BasicDaoMongo) Collection() *mongo.Collection {
	return bd.Mdb.Collection(bd.CollectionName())
}

func (bd *BasicDaoMongo) Disconnect() {
	ctx := context.Background()
	if err := bd.Mc.Disconnect(ctx); err != nil {
		panic(err)
	}
}
