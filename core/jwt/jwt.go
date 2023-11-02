package jwt

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	jv4 "github.com/golang-jwt/jwt/v4"
	"github.com/lucyanddarlin/lucy-ez-admin/config"
	rd "github.com/lucyanddarlin/lucy-ez-admin/core/redis"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

type jwt struct {
	redis *redis.Client
	conf  *config.JWT
	token string
}

type JWT interface {
	Compare(userID int64) bool
	IsExist(userId int64) bool
	Store(userID int64, token string, duration time.Duration) error
	Clear(userID int64) error
	Parse() (*types.Metadata, *jwtErr)
	Create(userID int64, data *types.Metadata) (string, error)
	IsWhiteList(method, path string) bool
	CheckUnique(userID int64) bool
}

func New(conf config.JWT, rd rd.Redis, token string) JWT {
	return &jwt{
		redis: rd.GetRedis(conf.Cache),
		conf:  &conf,
		token: token,
	}
}

func (j jwt) NewJwtErr(err string, opts ...jwtErrOption) *jwtErr {
	je := &jwtErr{
		err: errors.New(err),
	}
	for _, opt := range opts {
		opt(je)
	}
	return je
}

// uuid 获取存储的唯一 ID
func (j jwt) uuid(userID int64) string {
	return fmt.Sprintf("token_%x", j.encode(userID))
}

// encode 对存储数据进行加密编码
func (j jwt) encode(data any) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(data))))
}

// CheckUnique implements JWT.
func (j *jwt) CheckUnique(userID int64) bool {
	if !j.conf.Unique {
		return true
	}
	return j.Compare(userID)
}

// Clear implements JWT.
//
// 清楚用户缓存数据
func (j *jwt) Clear(userID int64) error {
	return j.redis.Del(context.Background(), j.uuid(userID)).Err()
}

// Compare implements JWT.
//
// 对比 token 是否和缓存一致
func (j *jwt) Compare(userID int64) bool {
	st, err := j.redis.Get(context.Background(), j.uuid(userID)).Result()
	if err != nil {
		return false
	}
	return st == j.encode(j.token)
}

// Create implements JWT.
//
// 生成 token 并缓存
func (j *jwt) Create(userID int64, data *types.Metadata) (string, error) {
	claims := make(jv4.MapClaims)
	claims["exp"] = time.Now().Unix() + int64(j.conf.Expire.Seconds())
	claims["iat"] = time.Now().Unix()
	claims["data"] = data
	tokenJwt := jv4.New(jv4.SigningMethodES256)
	tokenJwt.Claims = claims
	token, err := tokenJwt.SignedString([]byte(j.conf.Secret))
	if err != nil {
		return "", err
	}
	return token, j.Store(userID, token, j.conf.Expire)
}

// IsExist implements JWT.
//
// 判断缓存中是否存在用户的 token
func (j jwt) IsExist(userId int64) bool {
	st, _ := j.redis.Exists(context.Background(), j.uuid(userId)).Result()
	return st != 0
}

// IsWhiteList implements JWT.
func (j jwt) IsWhiteList(method string, path string) bool {
	return j.conf.WhiteList[strings.ToLower(method+":"+path)]
}

// Parse implements JWT.
//
// 解析用户的 token 信息
func (j jwt) Parse() (*types.Metadata, *jwtErr) {
	var m jv4.MapClaims = make(map[string]any)
	parser, err := jv4.ParseWithClaims(j.token, &m, func(token *jv4.Token) (interface{}, error) {
		return []byte(j.conf.Secret), nil
	})

	// 判断是否验证失败
	if err != nil && parser == nil {
		return nil, j.NewJwtErr("token verify error")
	}

	exp := int64(0)
	if m["exp"] != nil {
		exp = int64(m["exp"].(float64))
	}
	data := types.Metadata{}
	b, _ := json.Marshal(m["data"])
	_ = json.Unmarshal(b, &data)

	if err != nil {
		return &data, j.NewJwtErr(err.Error(),
			withVerify(parser.Valid),
			withExpired(errors.Is(err, jv4.ErrTokenExpired)),
			withExpiredUnix(exp),
			withRenewalUnix(int64(j.conf.Renewal.Seconds())),
		)
	}

	// 成功返回
	return &data, nil

}

// Store implements JWT.
//
// 存储用户 token 信息到存储中
func (j jwt) Store(userID int64, token string, duration time.Duration) error {
	return j.redis.Set(context.Background(), j.uuid(userID), j.encode(token), duration).Err()
}
