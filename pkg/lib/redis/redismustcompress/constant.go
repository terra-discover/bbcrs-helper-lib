package redismustcompress

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	viperRedisCompressionKey string = "REDIS_COMPRESSION"
	redisMustCompressKey     string = "must_compress"
)

type MustCompressValue string

const (
	mustCompressYes MustCompressValue = "yes"
	mustCompressNo  MustCompressValue = "no"
)

func NewMustCompressValue(mustCompress bool) (result MustCompressValue) {
	if !mustCompress {
		result = mustCompressNo
		return
	}

	result = mustCompressYes
	return
}

func ConvertStringToMustCompressValue(input string) (mustCompressValue MustCompressValue, err error) {
	switch input {
	case mustCompressYes.String():
		{
			mustCompressValue = mustCompressYes
			break
		}
	case mustCompressNo.String():
		{
			mustCompressValue = mustCompressNo
			break
		}
	default:
		{
			err = fmt.Errorf("redismustcompress.ConvertStringToMustCompressValue(): cannot convert %s to must compress value", input)
			break
		}
	}

	return
}

func GetEnvRedisCompression() bool {
	return viper.GetBool("REDIS_COMPRESSION")
}

func isRedisKeyNotFound(err error) (isNotFound bool) {
	isNotFound = err != nil && err == redis.Nil
	return
}

func (m MustCompressValue) IsTrue() bool {
	return m == mustCompressYes
}

func (m MustCompressValue) String() string {
	return string(m)
}

// On every init
// 1. Get mustCompress from redis
// 2. Get mustCompress from env
// 3. Compare mustCompress from redis and from env
//
// 3.a.mustCompress from redis is not found
// 4.a.save mustCompress from env to redis
// 5.a.set state mustCompress
//
// 3.b.mustCompress from redis == mustCompress from env
// 4.b.set state mustCompress
//
// 3.c.mustCompress from redis != mustCompress from env
// 4.c.got panic("REDIS_COMPRESSION environment has been changed to %s (,but in redis was %s). please flush all your cache before change the environment. you can also undo the changes.")
