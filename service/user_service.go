// Package service
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description:
 * @File:  user
 * @Version: 1.0.0
 * @Date: 2022/7/3 20:56
 */
package service

import (
	"errors"
	"r0Website-server/dao"
	"r0Website-server/global"
	"r0Website-server/middleware"
	"r0Website-server/models/bo"
	"r0Website-server/models/po"
	"r0Website-server/models/vo"
	"r0Website-server/utils"
	"time"
)

type UserService struct {
	UserDao *dao.UserDao `R0Ioc:"true"`
}

const UserColl = "users"

// UserLogin 用户登录
func (u *UserService) UserLogin(params vo.LoginVo) (ans *vo.LoginResultVo, err error) {
	var res po.User
	var result vo.LoginResultVo
	res, err = u.UserDao.FindObjByEmail(params.Email)
	if err != nil {
		return nil, errors.New("UserLogin 用户email不存在")
	}
	if !utils.IsPasswordMatch(params.Password, res.Password) {
		return nil, errors.New("UserLogin 用户密码错误")
	}
	token, err := middleware.GenToken(res)
	if err != nil {
		return nil, errors.New("UserLogin 构造Token失败，请联系管理员" + global.Config.Author.Email)
	}
	result.Token = token
	return &result, err
}

// UserRegister 用户注册
func (u *UserService) UserRegister(params vo.RegisterVo) (ans *vo.RegisterResultVo, err error) {
	var input po.User
	var result vo.RegisterResultVo
	updateUserInputByParams(&input, params)
	if emailCount := u.UserDao.EmailCount(input.Email); emailCount > 0 {
		return &result, &bo.UniqueError{UniqueField: "email", Msg: input.Email, Count: emailCount}
	}
	err = u.UserDao.CreateUser(input)
	result.Username = input.Username
	return &result, err
}

// updateUserInputByParams 根据参数更新需要输入的用户模型，同时加密密码
func updateUserInputByParams(input *po.User, params vo.RegisterVo) {
	input.Username = params.Username
	input.Salt, input.Password = utils.Encrypt(params.Password)
	input.Email = params.Email
	input.Phone = params.Phone
	curTime := time.Now()
	input.UpdateTime = curTime
	input.CreateTime = curTime
}
