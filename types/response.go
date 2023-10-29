package types

type Response struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Data     interface{} `json:"data"`
	TranceID string      `json:"tranceID"`
}

type ResponseList struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Data     interface{} `json:"data,omitempty"`
	TranceID string      `json:"tranceID"`
	Total    int64       `json:"total"`
}

func (r *Response) Error() string {
	return r.Msg
}
