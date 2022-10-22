// Package po
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 用户模型
 * @File:  user
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:49
 */
package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User 实体对应的数据表
// FOLLOW: https://www.mongodb.com/docs/drivers/go/current/usage-examples/struct-tagging/
type User struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"` // Mongo 主键 _id
	Username   string             `bson:"username"`      // 用户名
	Password   string             `bson:"password"`      // 密码
	UserLevel  int64              `bson:"user_level"`    // 用户水平
	IsLock     bool               `bson:"is_lock"`       // 是否锁定
	Email      string             `bson:"email"`         // 邮箱
	Phone      string             `bson:"phone"`         // 手机号
	NewTime    time.Time          `bson:"new_time"`      // 最近登陆时间
	CreateTime time.Time          `bson:"create_time"`   // 创建时间
	UpdateTime time.Time          `bson:"update_time"`   // 更新时间
	Brief      string             `bson:"brief"`         // 备注
	Salt       string             `bson:"salt"`          // 盐值
}
