// Package vo
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 登陆模型
 * @File:  login
 * @Version: 1.0.0
 * @Date: 2022/7/3 20:57
 */
package vo

// LoginVo 登录实体
type LoginVo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResultVo 登录返回数据
type LoginResultVo struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Brief    string `json:"brief" bson:"brief"`
	Token    string `json:"token" bson:"token"`
}
