package cert

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/lucyanddarlin/lucy-ez-admin/config"
)

type cert struct {
	mu sync.RWMutex
	m  map[string][]byte
}

type Cert interface {
	Get(name string) ([]byte, error)
	GetCert(name string) []byte
}

func New(cs []config.Cert) Cert {
	ct := cert{
		m:  make(map[string][]byte),
		mu: sync.RWMutex{},
	}

	ct.mu.Lock()
	defer ct.mu.Unlock()

	for _, item := range cs {
		file, err := os.Open(item.Path)
		if err != nil {
			panic("cert 初始化失败" + err.Error())
		}
		val, err := io.ReadAll(file)
		if err != nil {
			panic("读取 cert 证书失败" + err.Error())
		}
		ct.m[item.Name] = val
	}

	return &ct
}

// Get implements Cert.
//
//	@Description: 获取指定名称的证书, 不存在则返回报错
//	@param name 指定的证书名称
//	@return []byte
//	@return error
func (c *cert) Get(name string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.m[name] == nil {
		return nil, errors.New("mo exist cert")
	}
	return c.m[name], nil
}

// GetCert implements Cert.
//
//	@Description: 获取指定名称的证书, 不存在则返回 nil
//	@param name 指定的证书名称
//	@return []byte
//	@return error
func (c *cert) GetCert(name string) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.m[name]
}
