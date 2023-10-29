package http

import (
	"encoding/json"
	"time"
	"unsafe"

	"github.com/go-resty/resty/v2"
	"github.com/lucyanddarlin/lucy-ez-admin/config"
	"go.uber.org/zap"
)

type request struct {
	c        *config.Http
	request  *resty.Request
	logger   *zap.Logger
	inputLog bool
}

type response struct {
	err      error
	response *resty.Response
}

type RequestFunc func(*resty.Request) *resty.Request

type Request interface {
	DisableLog() Request
	Option(fn RequestFunc) Request
	Get(url string) (*response, error)
	Post(url string, data interface{}) (*response, error)
	PostJson(url string, data interface{}) (*response, error)
	Put(url string, data interface{}) (*response, error)
	PutJson(url string, data interface{}) (*response, error)
	Delete(url string) (*response, error)
}

func New(conf *config.Http, logger *zap.Logger) Request {
	client := resty.New()
	if conf.MaxRetryWaitTime == 0 {
		conf.RetryWaitTime = 5 * time.Second
	}
	if conf.Timeout == 0 {
		conf.Timeout = 60 * time.Second
	}
	client.SetRetryWaitTime(conf.RetryWaitTime)
	client.SetRetryMaxWaitTime(conf.MaxRetryWaitTime)
	client.SetRetryCount(conf.RetryCount)
	client.SetTimeout(conf.Timeout)
	return &request{
		request:  client.R(),
		logger:   logger,
		inputLog: true,
	}
}

func (r *request) log(t int64, res *response) {
	if !(r.c.EnableLog && r.inputLog) {
		return
	}

	resData := res.Body()
	logs := []zap.Field{
		zap.Any("method", r.request.Header),
		zap.Any("url", r.request.URL),
		zap.Any("header", r.request.Header),
		zap.Any("body", r.request.Body),
		zap.Any("cost", time.Now().UnixMilli()-t),
		zap.Any("res", *(*string)(unsafe.Pointer(&resData))),
	}
	if len(r.request.FormData) != 0 {
		logs = append(logs, zap.Any("form-data", r.request.FormData))
	}
	if len(r.request.QueryParam) != 0 {
		logs = append(logs, zap.Any("query", r.request.QueryParam))
	}
	r.logger.Info("request", logs...)
}

// DisableLog implements Request.
func (r *request) DisableLog() Request {
	r.inputLog = false
	return r
}

// Option implements Request.
func (r *request) Option(fn RequestFunc) Request {
	r.request = fn(r.request)
	return r
}

// Get implements Request.
func (r *request) Get(url string) (*response, error) {
	res := &response{}
	defer r.log(time.Now().UnixMilli(), res)
	res.response, res.err = r.request.Get(url)
	return res, res.err
}

// Post implements Request.
func (r *request) Post(url string, data interface{}) (*response, error) {
	res := &response{}
	defer r.log(time.Now().UnixMilli(), res)
	res.response, res.err = r.request.SetBody(data).Post(url)
	return res, res.err
}

// PostJson implements Request.
func (r *request) PostJson(url string, data interface{}) (*response, error) {
	res := &response{}
	defer r.log(time.Now().UnixMilli(), res)
	res.response, res.err = r.request.ForceContentType("application/json").SetBody(data).Post(url)
	return res, res.err
}

// Put implements Request.
func (r *request) Put(url string, data interface{}) (*response, error) {
	res := &response{}
	defer r.log(time.Now().UnixMilli(), res)
	res.response, res.err = r.request.SetBody(data).Put(url)
	return res, res.err
}

// PutJson implements Request.
func (r *request) PutJson(url string, data interface{}) (*response, error) {
	res := &response{}
	defer r.log(time.Now().UnixMilli(), res)
	res.response, res.err = r.request.ForceContentType("application/json").SetBody(data).Put(url)
	return res, res.err
}

// Delete implements Request.
func (r *request) Delete(url string) (*response, error) {
	res := &response{}
	defer r.log(time.Now().UnixMilli(), res)
	res.response, res.err = r.request.Delete(url)
	return res, res.err
}

func (r *response) Body() []byte {
	return r.response.Body()
}

func (r *response) Result(val interface{}) error {
	return json.Unmarshal(r.response.Body(), val)
}
