package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
	math_rand "math/rand"
	"time"
)

const (
	Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits        = "0123456789"
	LettersDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	HexChars      = "0123456789abcdef"
)

func init() {
	// 初始化随机数种子
	math_rand.Seed(time.Now().UnixNano())
}

// String 生成指定长度的随机字符串
// length: 字符串长度
// charset: 字符集，可以使用预定义的常量，如Letters, Digits, LettersDigits
func String(length int, charset string) string {
	if charset == "" {
		charset = LettersDigits
	}

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果加密随机数生成失败，回退到math/rand
			b[i] = charset[math_rand.Intn(len(charset))]
		} else {
			b[i] = charset[n.Int64()]
		}
	}
	return string(b)
}

// RandomLetters 生成指定长度的随机字母字符串
func RandomLetters(length int) string {
	return String(length, Letters)
}

// RandomDigits 生成指定长度的随机数字字符串
func RandomDigits(length int) string {
	return String(length, Digits)
}

// RandomLettersDigits 生成指定长度的随机字母数字字符串
func RandomLettersDigits(length int) string {
	return String(length, LettersDigits)
}

// Hex 生成指定长度的随机十六进制字符串
func Hex(length int) string {
	return String(length, HexChars)
}

// Int 生成指定范围内的随机整数
// min: 最小值（包含）
// max: 最大值（包含）
func Int(min, max int) int {
	if min >= max {
		return min
	}
	return min + math_rand.Intn(max-min+1)
}

// UUID 生成UUID v4
func UUID() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		// 如果加密随机数生成失败，回退到math/rand
		for i := range u {
			u[i] = byte(math_rand.Intn(256))
		}
	}

	// 设置版本和变体
	u[6] = (u[6] & 0x0F) | 0x40 // 版本 4
	u[8] = (u[8] & 0x3F) | 0x80 // 变体 RFC4122

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// GenerateRandomString 生成指定长度的随机字母数字字符串
// 这是RandomLettersDigits的别名，用于生成应用密钥
// length: 字符串长度
func GenerateRandomString(length int) string {
	return RandomLettersDigits(length)
}
