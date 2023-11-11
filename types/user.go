package types

type UserLoginRequest struct {
	Phone       string `json:"phone"  binding:"required"`
	Password    string `json:"password"  binding:"required"`
	CaptchaName string `json:"-"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	Captcha     string `json:"captcha"  binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type Password struct {
	Password string `json:"password"`
	Time     int64  `json:"time"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"id" binding:"required"`
	TeamID   int64  `json:"team_id"`
	RoleID   int64  `json:"role_Id"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Sex      *bool  `json:"sex"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   *bool  `json:"status"`
	Password string `json:"password"`
}

type AddUserRequest struct {
	TeamID   int64  `json:"team_id" binding:"required"`
	RoleID   int64  `json:"role_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Sex      *bool  `json:"sex" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" binding:"required"`
	Status   *bool  `json:"status" binding:"required"`
}
