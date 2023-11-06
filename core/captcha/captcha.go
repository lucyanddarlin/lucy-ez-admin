package captcha

import (
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"

	"github.com/lucyanddarlin/lucy-ez-admin/config"
	e "github.com/lucyanddarlin/lucy-ez-admin/core/email"
	"github.com/lucyanddarlin/lucy-ez-admin/core/redis"
)

type res struct {
	ID     string `json:"id"`
	Base64 string `json:"base64,omitempty"`
	Expire int    `json:"expire"`
}

type captcha struct {
	mu    sync.RWMutex
	cache redis.Redis
	email e.Email
	m     map[string]config.Captcha
}

type Captcha interface {
	Image(ip, name string) Image
	Email(ip, name string) Email
}

// New 初始化 captcha 实例
func New(cs []config.Captcha, rs redis.Redis, email e.Email) Captcha {
	cpIns := captcha{
		cache: rs,
		email: email,
		m:     make(map[string]config.Captcha),
		mu:    sync.RWMutex{},
	}

	cpIns.mu.Lock()
	defer cpIns.mu.Unlock()

	// 配置 Config.Captcha 中的每个配置项存储至 m 中
	for _, item := range cs {
		cpIns.m[item.Name+":"+item.Type] = item
	}

	return &cpIns
}

// getTemplate 获取指定模板
func (c *captcha) getTemplate(name, tp string) (config.Captcha, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	template, is := c.m[name+":"+tp]
	return template, is
}

// cid 获取用户对应验证码场景的唯一 id
func (c *captcha) cid(ip, name, tp string) string {
	return fmt.Sprintf("captcha:%s:%s:%x", name, tp, md5.Sum([]byte(ip)))
}

// randomCode 随机生成验证码
func (c *captcha) randomCode(len int) string {
	r := rand.New(rand.NewSource(int64(len)))
	var code = r.Intn(int(math.Pow10(len)) - int(math.Pow10(len-1)))
	return strconv.Itoa(code + int(math.Pow10(len-1)))
}

// Image implements Captcha.
//
// 实例化图形验证码
func (c *captcha) Image(ip string, name string) Image {
	return &image{
		name:    name,
		tp:      "image",
		captcha: c,
		ip:      ip,
	}
}

// Email implements Captcha.
//
// 实例化邮箱验证码
func (c *captcha) Email(ip string, name string) Email {
	return &email{
		name:    name,
		tp:      "email",
		captcha: c,
		ip:      ip,
	}
}
