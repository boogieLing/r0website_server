// Package utils
/**
 * @Author: r0
 * @Mail: boogieLing_o@qq.com
 * @Description: 加密工具 FOLLOW FreeBSD Colin Percival
 * @File:  encrypt
 * @Version: 1.0.0
 * @Date: 2022/7/4 16:18
 */
package utils

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/sony/sonyflake"
	"golang.org/x/crypto/scrypt"
	"log"
	"strconv"
)

func Encrypt(password string) (string, string) {
	has := md5.Sum([]byte(password))
	salt := has[:8]
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(dk)
}

func IsPasswordMatch(rawPassword, dbPassword string) bool {
	_, crypt := Encrypt(rawPassword)
	return crypt == dbPassword
}
func GenSonyflake() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	// Note: this is base16, could shorten by encoding as base62 string
	// fmt.Printf("github.com/sony/sonyflake:   %x\n", id)
	return strconv.Itoa(int(id))
}
