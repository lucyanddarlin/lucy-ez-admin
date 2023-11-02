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
