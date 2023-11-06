package captcha

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type email struct {
	name    string
	tp      string
	captcha *captcha
	ip      string
}

type Email interface {
	// New 发送邮件验证码
	//
	//  @param email 发送邮箱
	//  @return res 验证码 id, 过期时间
	//  @return error 验证码通过返回 nil, 否则返回定义错误
	New(email string) (*res, error)
	// Verify
	//
	//  @param id 验证码 id
	//  @param answer 用户提交的验证码
	//  @return 验证通过返回 nil, 否则返回定义错误
	Verify(id, answer string) error
}

func (e *email) New(email string) (*res, error) {
	// 获取指定模板的配置
	cp, is := e.captcha.getTemplate(e.name, e.tp)
	if !is {
		return nil, errors.New(fmt.Sprintf("%s captcha is not exist", e.tp))
	}

	// 生成随机验证码
	answer := e.captcha.randomCode(cp.Length)

	// 获取验证码存储器
	cache := e.captcha.cache.GetRedis(cp.Cache)

	// 获取当前用户的场景唯一 id
	cid := e.captcha.cid(e.ip, e.name, e.tp)

	// 清除上一次生成的结果, 防止同时造成大量生成请求占用内存
	if id, _ := cache.Get(context.Background(), cid).Result(); id != "" {
		cache.Del(context.Background(), id)
	}

	// 获取当前验证码唯一 id
	uid := uuid.New().String()
	if err := cache.Set(context.Background(), uid, answer, cp.Expire).Err(); err != nil {
		return nil, err
	}

	// 将本次验证码挂载到当前场景 id 上
	if err := cache.Set(context.Background(), cid, uid, cp.Expire).Err(); err != nil {
		return nil, err
	}

	// 发送邮件
	if err := e.captcha.email.NewSender(cp.Template).Send(email, map[string]any{
		"answer": answer,
		"minute": int(cp.Expire.Minutes()),
	}); err != nil {
		return nil, err
	}

	// 返回生成结果
	return &res{
		ID:     uid,
		Expire: int(cp.Expire.Seconds()),
	}, nil
}

func (e *email) Verify(id, answer string) error {
	return nil
}
