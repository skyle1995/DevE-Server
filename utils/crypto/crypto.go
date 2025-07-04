package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// MD5 计算字符串的MD5哈希值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 计算字符串的SHA1哈希值
func SHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 计算字符串的SHA256哈希值
func SHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SignParams 对参数进行签名
// params: 参数map
// secret: 密钥
// 签名规则：
// 1. 将参数按照参数名的字典序排序
// 2. 将参数名和参数值拼接成字符串，格式为：参数名=参数值
// 3. 将拼接后的字符串用&连接，再加上密钥
// 4. 对最终的字符串进行MD5加密，得到签名
func SignParams(params map[string]string, secret string) string {
	// 获取所有参数名
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	// 按字典序排序
	sort.Strings(keys)

	// 拼接参数
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	// 拼接密钥
	signStr := strings.Join(parts, "&") + secret

	// 计算MD5
	return MD5(signStr)
}

// VerifySign 验证签名
// params: 参数map
// sign: 签名
// secret: 密钥
func VerifySign(params map[string]string, sign string, secret string) bool {
	// 计算签名
	calcSign := SignParams(params, secret)

	// 比较签名
	return calcSign == sign
}
