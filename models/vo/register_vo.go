// Package vo
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description:
 * @File:  register
 * @Version: 1.0.0
 * @Date: 2022/7/4 15:06
 */
package vo

// RegisterVo 注册实体
type RegisterVo struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
}

// RegisterResultVo 登录返回数据
type RegisterResultVo struct {
	Username string `json:"username" bson:"username"`
}
