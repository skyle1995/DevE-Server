package captcha

import (
	"bytes"
	"github.com/dchest/captcha"
	"net/http"
	"time"
)

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Width      int           // 验证码图片宽度
	Height     int           // 验证码图片高度
	Length     int           // 验证码长度
	Expiration time.Duration // 验证码过期时间
}

// DefaultConfig 默认验证码配置
func DefaultConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Width:      160,
		Height:     80,
		Length:     4,
		Expiration: 10 * time.Minute,
	}
}

// Captcha 验证码管理器
type Captcha struct {
	config *CaptchaConfig
}

// New 创建一个新的验证码管理器
func New(config *CaptchaConfig) *Captcha {
	if config == nil {
		config = DefaultConfig()
	}

	// 设置验证码过期时间
	captcha.SetCustomStore(captcha.NewMemoryStore(1000, config.Expiration))

	return &Captcha{
		config: config,
	}
}

// Generate 生成一个新的验证码ID
func (c *Captcha) Generate() string {
	return captcha.NewLen(c.config.Length)
}

// Verify 验证验证码是否正确
func (c *Captcha) Verify(id, value string) bool {
	return captcha.VerifyString(id, value)
}

// WriteImage 将验证码图片写入HTTP响应
func (c *Captcha) WriteImage(w http.ResponseWriter, id string) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")

	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, id, c.config.Width, c.config.Height); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}