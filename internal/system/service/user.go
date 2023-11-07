package service

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/forgoer/openssl"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

const (
	encodePasswordCert  = "encodePassword"
	decodePasswordCert  = "decodePassword"
	passwordExpiredTime = 10 * 1000
)

// UserLogin 用户登录
func UserLogin(ctx *core.Context, in *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	resp = new(types.UserLoginResponse)
	// 函数返回时的错误处理
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
	// TODO: 后期开启
	// if time.Now().UnixMilli()-pw.Time > passwordExpiredTime {
	// 	err = errors.PasswordExpireError
	// 	return
	// }

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
	if !role.RoleStatus(ctx, user.RoleID) {
		err = errors.RoleDisableError
		return
	}

	// 对比用户密码
	if !tools.CompareHashPwd(password, in.Password) {
		err = errors.PasswordError
		return
	}

	// 生成登陆的 token
	if resp.Token, err = ctx.Jwt().Create(user.ID, &types.Metadata{
		UserID:    user.ID,
		RoleID:    user.RoleID,
		RoleKey:   user.Role.Keyword,
		DataScope: user.Role.DataScope,
		Username:  user.Name,
		TeamID:    user.TeamID,
	}); err != nil {
		return nil, err
	}

	// // 修改登录时间
	return resp, user.UpdateLastLogin(ctx, time.Now().Unix())

}

// UserLogout 用户退出登录
func UserLogout(ctx *core.Context) error {
	metadata, _ := ctx.Jwt().Parse()
	if metadata != nil {
		return ctx.Jwt().Clear(metadata.UserID)
	}
	return nil
}

// RefreshToken 用户刷新 token
func RefreshToken(ctx *core.Context) (*types.UserLoginResponse, error) {
	md, err := ctx.Jwt().Parse()
	if md == nil {
		return nil, errors.MetadataError
	}
	if err == nil {
		return nil, errors.RefreshActiveTokenError
	}
	if !err.CanRenewal() {
		return nil, errors.RefTokenExpiredError
	}

	token, e := ctx.Jwt().Create(md.UserID, md)
	if e != nil {
		return nil, e
	}

	return &types.UserLoginResponse{
		Token: token,
	}, e
}

// CurrentUser 获取当前登录用户信息
func CurrentUser(ctx *core.Context) (*model.User, error) {
	md := ctx.Metadata()
	if md == nil {
		return nil, errors.MetadataError
	}

	user := model.User{}
	return &user, user.OneByID(ctx, md.UserID)
}
