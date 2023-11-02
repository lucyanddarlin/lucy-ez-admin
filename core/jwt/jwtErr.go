package jwt

import "time"

type jwtErr struct {
	err       error
	isVerify  bool
	isExpired bool
	exp       int64
	renewal   int64
}

type jwtErrOption func(err *jwtErr)

func withVerify(val bool) jwtErrOption {
	return func(err *jwtErr) {
		err.isVerify = val
	}
}

func withExpired(val bool) jwtErrOption {
	return func(err *jwtErr) {
		err.isExpired = val
	}
}

func withExpiredUnix(val int64) jwtErrOption {
	return func(err *jwtErr) {
		err.exp = val
	}
}

func withRenewalUnix(val int64) jwtErrOption {
	return func(err *jwtErr) {
		err.renewal = val
	}
}

func (j *jwtErr) Error() string {
	return j.err.Error()
}

// IsVerify 是否验证通过
func (j *jwtErr) IsVerify() bool {
	return j.isVerify
}

// isExpired 是否过期
func (j *jwtErr) IsExpired() bool {
	return j.isExpired
}

// CanRenewal 是否能够续期
func (j *jwtErr) CanRenewal() bool {
	return time.Now().Unix()-j.exp <= j.renewal
}
