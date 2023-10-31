package types

type UserLoginRequest struct {
	Phone       string `json:"phone"  binding:"required"`
	Password    string `json:"password"  binding:"required"`
	CaptchaName string `json:"-"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	Captcha     string `json:"captcha"  binding:"required"`
}
