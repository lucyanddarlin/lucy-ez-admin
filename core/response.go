package core

import "github.com/lucyanddarlin/lucy-ez-admin/types"

// RespSuccess
//
//	@Description: 返回成功
func (ctx *Context) RespSuccess() {
	ctx.JSON(200, &types.Response{
		Code:     ctx.Config().Service.SuccessCode,
		Msg:      "success",
		TranceID: ctx.TranceID(),
	})
}

// RespData
//
//	@Description: 返回成功并且携带数据
//	@param data 成功返回的数据
func (c *Context) RespData(data any) {
	c.JSON(200, &types.Response{
		Code:     c.Config().Service.SuccessCode,
		Msg:      "success",
		Data:     data,
		TranceID: c.TranceID(),
	})
}

// RespList
//
//	@Description: 返回成功并且携带列表数量,用于分页查询
//	@param total 总的数量条数
//	@param data 分页查询的数据
func (c *Context) RespList(total int64, data any) {
	c.JSON(200, &types.ResponseList{
		Code:     c.Config().Service.SuccessCode,
		Msg:      "success",
		Data:     data,
		Total:    total,
		TranceID: c.TranceID(),
	})
}

// RespError
//
//		@Description: 返回数据错误的信息
//	 @param err
func (c *Context) RespError(err error) {
	if response, is := err.(*types.Response); is {
		response.TranceID = c.TranceID()
		c.JSON(200, response)
	} else {
		c.JSON(200, &types.Response{
			Code:     c.Config().Service.ErrorCode,
			Msg:      err.Error(),
			TranceID: c.TranceID(),
		})
	}
}
