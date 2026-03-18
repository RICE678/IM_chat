package dao

import (
	"math/rand"
	"regexp"
	"strconv"
)

// VerifyEmailFormat 验证邮箱格式
func VerifyEmailFormat(email string) bool {
	pattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// GetConfirmCode 获取六位数随机验证码函数
func GetConfirmCode() string {
	var confirmCode int
	for i := 0; i < 6; i++ {
		confirmCode = confirmCode*10 + (rand.Intn(9) + 1)

	}
	confirmCodeStr := strconv.Itoa(confirmCode)
	return confirmCodeStr
}
