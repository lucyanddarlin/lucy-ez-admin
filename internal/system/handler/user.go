package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	// 创建参数 struct, 定义 captchaName 对应模板 login
	in := types.UserLoginRequest{
		CaptchaName: "login",
	}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if resp, err := service.UserLogin(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if err := service.UserLogout(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// RefreshToken 刷新 token
func RefreshToken(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if resp, err := service.RefreshToken(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

// CurUser 获取当前用户信息
func CurUser(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if user, err := service.CurrentUser(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(user)
	}
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UpdateUserRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.UpdateUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// AddUser 添加用户
func AddUser(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.AddUserRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.AddUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// UpdateUserInfo 更新当前用户信息
func UpdateUserInfo(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UpdateUserInfoRequest{}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.RespError(err)
		return
	}

	if err := service.UpdateCurrentUserInfo(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// UpdateUserInfoByVerify 更新用户重要信息
func UpdateUserInfoByVerify(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UpdateUserInfoByVerifyRequest{
		CaptchaName: "user",
	}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.UpdateUserInfoByVerify(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.DeleteUserRequest{}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.DeleteUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// PageUser 用户列表分页
func PageUser(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.PageUserRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if list, total, err := service.PageUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespList(total, list)
	}
}
