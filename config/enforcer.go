package config

type Enforcer struct {
	Enable    bool
	DB        string
	WhiteList map[string]bool
	JWT       *JWT
}
