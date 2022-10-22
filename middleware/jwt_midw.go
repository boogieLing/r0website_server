// Package middleware
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: Jwt中间件 校验、生成
 * @File:  Jwt
 * @Version: 1.0.0
 * @Date: 2022/7/3 18:47
 */
package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"r0Website-server/global"
	"r0Website-server/models/po"
	"strings"
	"time"
)

type WebsiteClaims struct {
	User po.User `json:"user"`
	jwt.StandardClaims
}

// ParseToken 解析token
func ParseToken(tokenString string) (*WebsiteClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &WebsiteClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.SignKey), nil
	})
	//if err != nil {
	//	return nil, err
	//}
	claims, ok := token.Claims.(*WebsiteClaims)
	if ok && token.Valid {
		return claims, nil
	} else if ok && !token.Valid {
		// 能解析说明格式正确，非法说明大概率是过期，应该返回claims
		return claims, err
	}
	return nil, errors.New("非法的Token串")
}

// Jwt jwt鉴权控制
// - 前置校验
// - Jwt 本身的校验
// - 过期5分钟内自动续期
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		// 先判空
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 判断token的前缀是否正确
		items := strings.SplitN(token, " ", 2)
		if !(len(items) == 2 && items[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// 解析token信息
		blogClaims, err := ParseToken(items[1])
		if err != nil {
			// 如果过期时间没超过5分钟，自动续期
			var boolRenewSuccess = false
			if strings.Contains(err.Error(), "expired") {
				newToken, _ := RenewToken(blogClaims)
				if newToken != "" {
					// 续期成功
					c.Header("new_token", newToken)
					c.Request.Header.Set("Authorization", newToken)
					boolRenewSuccess = true
				}
			}
			if !boolRenewSuccess {
				// 续期失败，还是原地abort吧
				c.JSON(http.StatusOK, gin.H{
					"code": 2005,
					"msg":  "无效的Token: " + err.Error(),
				})
				c.Abort()
				return
			}
		}
		// 将当前请求的user保存到请求的上下文c上
		c.Set("userInfo", blogClaims.User)
		c.Next()
		// 后续的处理函数可以用过c.Get("userInfo")来获取当前请求的用户信息
		// 类型是any 需要类型推导
	}
}

// GenToken 构建Token
// - 构建一个声明
// - 创建签名对象
// - 签发编码字符串 sign-key在config/config.yml
func GenToken(u po.User) (string, error) {
	c := WebsiteClaims{
		u,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                 // 生效时间
			ExpiresAt: time.Now().Unix() + global.Config.JWT.ExpiresTime, // 过期时间
			Issuer:    global.Config.JWT.Issuer,                          // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(global.Config.JWT.SignKey))
}
func RenewToken(claims *WebsiteClaims) (string, error) {
	// 最多允许过期5分钟
	if withinLimit(claims.ExpiresAt, 3*100) {
		return GenToken(claims.User)
	}
	return "", errors.New("用户: " + claims.User.Username + " 已过期")
}

// 计算时间是否超过限制
func withinLimit(s int64, limit int64) bool {
	return time.Now().Unix()-s < limit
}
