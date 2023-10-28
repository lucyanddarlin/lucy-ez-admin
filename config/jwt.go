package config

import "time"

type JWT struct {
	Enable    bool
	Cache     string
	Header    string
	Secret    string
	Expire    time.Duration
	Renewal   time.Duration
	Unique    bool
	WhiteList map[string]bool
}
