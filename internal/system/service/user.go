package service

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/forgoer/openssl"
	"github.com/jinzhu/copier"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"gorm.io/gorm"
)

const (
	encodePasswordCert  = "encodePassword"
	decodePasswordCert  = "decodePassword"
	passwordExpiredTime = 10 * 1000
)

// CurrentAdminTeamIds 获取当前用户管理的部门 id
func CurrentAdminTeamIds(ctx *core.Context) ([]int64, error) {
	md := ctx.Metadata()
	if md == nil {
		return nil, errors.MetadataError
	}

	user := model.User{}
	ids, err := user.GetAdminTeamIdByUserId(ctx, md.UserID)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// CurrentAdminRoleIds 获取当前用户管理的角色 ids
func CurrentAdminRoleIds(ctx *core.Context) ([]int64, error) {
	md := ctx.Metadata()
	if md == nil {
		return nil, errors.MetadataError
	}

	role := model.Role{}
	roleTree, err := role.Tree(ctx, md.RoleID)
	if err != nil {
		return nil, err
	}
	return tree.GetTreeID(roleTree), nil
}

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

// UpdateUser 更新用户信息
func UpdateUser(ctx *core.Context, in *types.UpdateUserRequest) error {
	user := model.User{}
	if user.OneByID(ctx, in.ID) != nil {
		return errors.DBNotFoundError
	}

	// 超级管理员不允许修改所在部门和角色
	if in.ID == 1 {
		in.RoleID = 0
		in.TeamID = 0
		if *user.Status != *in.Status {
			return errors.SuperAdminDelError
		}
	}

	// 修改角色时, 也只允许修改到自己所在管辖的角色
	if in.RoleID != 0 && in.RoleID != user.RoleID {
		ids, err := CurrentAdminTeamIds(ctx)
		if err != nil {
			return err
		}
		if !tools.InList(ids, in.RoleID) {
			return errors.NotEditUserRoleError
		}
	}

	// 获取用户可以管理的部门
	ids, err := CurrentAdminTeamIds(ctx)
	if err != nil {
		return nil
	}
	// 只允许更新当前部门的用户信息
	if !tools.InList(ids, user.TeamID) {
		return errors.NotEditTeamError
	}

	// 修改部门时, 只允许修改到自己所管辖的部门
	if in.TeamID != 0 && in.TeamID != user.TeamID && !tools.InList(ids, in.TeamID) {
		return errors.NotEditTeamError
	}

	if copier.Copy(&user, in) != nil {
		return errors.AssignError
	}
	return user.Update(ctx)
}

// AddUser 添加用户
func AddUser(ctx *core.Context, in *types.AddUserRequest) error {
	user := model.User{}
	user.OneByName(ctx, in.Name)
	if user.Name != "" {
		return errors.ExistUserNameError
	}

	user = model.User{}

	if in.Nickname == "" {
		in.Nickname = in.Name
	}

	if err := copier.Copy(&user, in); err != nil {
		return err
	}

	// 获取操作者所管理的部门
	ids, err := CurrentAdminTeamIds(ctx)
	if err != nil {
		return err
	}

	// 添加用户时, 只允许添加到当前操作者所管理的部门
	if !tools.InList(ids, in.TeamID) {
		return errors.NotAddTeamUserError
	}

	return user.Create(ctx)
}

// UpdateCurrentUserInfo 更新用户信息
func UpdateCurrentUserInfo(ctx *core.Context, in *types.UpdateUserInfoRequest) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	user := model.User{}
	if err := copier.Copy(&user, in); err != nil {
		return err
	}
	user.ID = md.UserID

	return user.Update(ctx)
}

// DeleteUser 删除用户
func DeleteUser(ctx *core.Context, in *types.DeleteUserRequest) error {
	if in.ID == 1 {
		return errors.SuperAdminDelError
	}

	user := model.User{}
	if user.OneByID(ctx, in.ID) != nil {
		return errors.DBNotFoundError
	}

	// 获取操作者所管理的部门
	ids, err := CurrentAdminTeamIds(ctx)
	if err != nil {
		return err
	}

	// 只允许删除当前操作者所管理部门的人员
	if !tools.InList(ids, user.TeamID) {
		return errors.NotDelTeamUserError
	}
	return user.DeleteByID(ctx, in.ID)
}

// PageUser 获取用户列表分页
func PageUser(ctx *core.Context, in *types.PageUserRequest) ([]*model.User, int64, error) {
	user := model.User{}

	// 获取当前操作者所管理的部门
	ids, err := CurrentAdminTeamIds(ctx)
	if err != nil {
		return nil, 0, err
	}

	return user.Page(ctx, types.PageOptions{
		Page:     in.Page,
		PageSize: in.PageSize,
		Model:    in,
		Scopes: func(db *gorm.DB) *gorm.DB {
			return db.Where("team_id in ?", ids)
		},
	})

}
