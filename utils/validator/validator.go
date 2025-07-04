package validator

import (
	"regexp"
	"strings"
)

// IsEmail 验证是否为有效的电子邮件地址
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// IsPhone 验证是否为有效的手机号码（中国大陆）
func IsPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(phone)
}

// IsURL 验证是否为有效的URL
func IsURL(url string) bool {
	pattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(url)
}

// IsIP 验证是否为有效的IP地址
func IsIP(ip string) bool {
	pattern := `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(ip)
}

// IsIPv6 验证是否为有效的IPv6地址
func IsIPv6(ip string) bool {
	pattern := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(ip)
}

// IsIDCard 验证是否为有效的身份证号码（中国大陆）
func IsIDCard(idCard string) bool {
	// 18位身份证号码正则表达式
	pattern := `^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(idCard)
}

// IsNumeric 验证是否为数字
func IsNumeric(str string) bool {
	pattern := `^[0-9]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

// IsAlpha 验证是否为字母
func IsAlpha(str string) bool {
	pattern := `^[a-zA-Z]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

// IsAlphaNumeric 验证是否为字母或数字
func IsAlphaNumeric(str string) bool {
	pattern := `^[a-zA-Z0-9]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

// IsEmpty 验证是否为空
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNotEmpty 验证是否非空
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// IsLengthBetween 验证字符串长度是否在指定范围内
func IsLengthBetween(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

// IsStrongPassword 验证是否为强密码
// 强密码要求：
// 1. 长度至少为8位
// 2. 包含至少一个大写字母
// 3. 包含至少一个小写字母
// 4. 包含至少一个数字
// 5. 包含至少一个特殊字符
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasUppercase && hasLowercase && hasDigit && hasSpecial
}
