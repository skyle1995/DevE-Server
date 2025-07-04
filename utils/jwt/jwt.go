package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT声明结构体
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Config JWT配置结构体
type Config struct {
	SigningKey     string        // 签名密钥
	Issuer         string        // 签发者
	ExpiresTime    time.Duration // 过期时间
	RefreshTime    time.Duration // 刷新时间
	TokenHeadName  string        // Token头名称
	TokenLookup    string        // Token查找位置
	TokenBlacklist bool          // 是否启用Token黑名单
}

// DefaultConfig 默认JWT配置
func DefaultConfig() Config {
	return Config{
		SigningKey:     "default_signing_key",
		Issuer:         "deveSystem",
		ExpiresTime:    time.Hour * 24,
		RefreshTime:    time.Hour * 24 * 7,
		TokenHeadName:  "Bearer",
		TokenLookup:    "header:Authorization",
		TokenBlacklist: false,
	}
}

// JWT JWT管理器结构体
type JWT struct {
	Config Config
}

// New 创建新的JWT管理器
func New(config Config) *JWT {
	return &JWT{
		Config: config,
	}
}

// CreateToken 创建JWT Token
func (j *JWT) CreateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Config.Issuer,
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Config.ExpiresTime)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Config.SigningKey))
}

// ParseToken 解析JWT Token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Config.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新JWT Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查是否需要刷新
	if time.Until(claims.ExpiresAt.Time) > j.Config.ExpiresTime/2 {
		return tokenString, nil // Token还有足够的有效期，不需要刷新
	}

	// 创建新的Token
	return j.CreateToken(claims.UserID, claims.Username, claims.Role)
}

// IsTokenExpired 检查Token是否过期
func (j *JWT) IsTokenExpired(claims *Claims) bool {
	return claims.ExpiresAt.Time.Before(time.Now())
}

// GetTokenFromHeader 从HTTP请求头获取Token
func (j *JWT) GetTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("auth header is empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != j.Config.TokenHeadName {
		return "", errors.New("auth header format is wrong")
	}

	return parts[1], nil
}

// GetTokenFromCookie 从HTTP请求的Cookie中获取Token
func (j *JWT) GetTokenFromCookie(cookieValue string) (string, error) {
	if cookieValue == "" {
		return "", errors.New("cookie is empty")
	}

	return cookieValue, nil
}

// GetUserIDFromToken 从Token中获取用户ID
func (j *JWT) GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// GetUsernameFromToken 从Token中获取用户名
func (j *JWT) GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Username, nil
}

// GetRoleFromToken 从Token中获取角色
func (j *JWT) GetRoleFromToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Role, nil
}

// GetExpirationFromToken 从Token中获取过期时间
func (j *JWT) GetExpirationFromToken(tokenString string) (time.Time, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

// GetTokenRemainingTime 获取Token剩余有效时间
func (j *JWT) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	return time.Until(claims.ExpiresAt.Time), nil
}

// GetTokenFromHeader 从HTTP请求头获取Token的辅助函数
func GetTokenFromHeader(authHeader, tokenHeadName string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != tokenHeadName {
		return ""
	}

	return parts[1]
}
