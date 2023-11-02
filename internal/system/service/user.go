package service

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/forgoer/openssl"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

const (
	encodePasswordCert = "encodePassword"
	decodePasswordCert = "decodePassword"
)

// UserLogin 用户登录
func UserLogin(ctx *core.Context, in *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	resp = new(types.UserLoginResponse)
	defer func() {
		if !(errors.Is(err, errors.UserDisableError) || errors.Is(err, errors.CaptchaError)) {
			_ = AddLoginLog(ctx, in.Phone, err)
		}
	}()

	// 判断验证码是否正确
	if err = ctx.ImageCaptcha(in.CaptchaName).Verify(in.CaptchaID, in.Captcha); err != nil {
		err = errors.CaptchaError
		return
	}

	// 密码解密
	passByte, _ := base64.StdEncoding.DecodeString(in.Password)
	decryptData, err := openssl.RSADecrypt(passByte, ctx.Cert().GetCert(decodePasswordCert))
	if err != nil {
		err = errors.RsaPasswordError
		return
	}

	// 序列化密码
	var pw types.Password
	if json.Unmarshal(decryptData, &pw) != nil {
		err = errors.RsaPasswordError
		return
	}

	// 判断当前时间戳是否过期,超过 10s 则拒绝
	if time.Now().UnixMilli()-pw.Time > 10*1000 {
		err = errors.PasswordExpireError
		return
	}

	in.Password = pw.Password

	// 通过手机号获取信息
	user := model.User{}
	if err = user.OneByPhone(ctx, in.Phone); err != nil {
		err = errors.UserNotFoundError
		return
	}

	// 由于屏蔽了 password, 需要调用指定方法查询
	password, err := user.PasswordByPhone(ctx, in.Phone)
	if err != nil {
		err = errors.UserNotFoundError
		return
	}

	// 用户被禁用则拒绝登陆
	if !*user.Status {
		err = errors.UserDisableError
		return
	}

	// 所属角色被禁用则拒绝登录
	role := model.Role{}
	// if role.

}
