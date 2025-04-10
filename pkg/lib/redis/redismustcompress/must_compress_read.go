package redismustcompress

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type IRedisGet interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

type BuilderMustCompressRead struct {
	RedisGet IRedisGet
}

type MustCompressRead struct {
	builder            BuilderMustCompressRead
	mustCompressString string
	isFound            bool
	mustCompressValue  MustCompressValue
	mustCompress       bool
	err                error
}

func NewMustCompressRead(builder BuilderMustCompressRead) (mcs *MustCompressRead) {
	mcs = new(MustCompressRead)

	// builder
	mcs.setBuilder(builder)
	return
}

func (mcs *MustCompressRead) Read() {
	// clear err first
	mcs.clearErr()

	// get mustCompressString from redis
	mustCompressString, errString := mcs.getMustCompressStringFromRedis()
	if isRedisKeyNotFound(errString) {
		mcs.setIsFound(false)
		return
	}
	mcs.setIsFound(true)

	if errString != nil {
		mcs.setErr("getMustCompressStringFromRedis()", errString)
		return
	}
	mcs.setMustCompressString(mustCompressString)

	// mustCompressValue
	mustCompressVal, errVal := mcs.genMustCompressValue()
	if errVal != nil {
		mcs.setErr("genMustCompressValue()", errVal)
		return
	}
	mcs.setMustCompressValue(mustCompressVal)

	// mustCompress
	mustCompress := mcs.genMustCompress()
	mcs.setMustCompress(mustCompress)
}

func (mcs MustCompressRead) Me() MustCompressRead {
	return mcs
}

func (r MustCompressRead) IsFound() (isFound bool) {
	isFound = r.getIsFound()
	return
}

func (r MustCompressRead) MustCompress() (mustCompress bool) {
	mustCompress = r.getMustCompress()
	return
}

func (mcs MustCompressRead) Err() (err error) {
	err = mcs.getErr()
	return
}

func (r *MustCompressRead) setBuilder(builder BuilderMustCompressRead) {
	r.builder = builder
}

func (r MustCompressRead) getBuilder() (builder BuilderMustCompressRead) {
	builder = r.builder
	return
}

func (r *MustCompressRead) setIsFound(isFound bool) {
	r.isFound = isFound
}

func (r MustCompressRead) getIsFound() (isFound bool) {
	isFound = r.isFound
	return
}

func (r *MustCompressRead) setMustCompress(mustCompress bool) {
	r.mustCompress = mustCompress
}

func (r MustCompressRead) getMustCompress() (mustCompress bool) {
	mustCompress = r.mustCompress
	return
}

func (r MustCompressRead) genMustCompress() (mustCompress bool) {
	mustCompressVal := r.getMustCompressValue()
	mustCompress = mustCompressVal.IsTrue()
	return
}

func (r *MustCompressRead) setMustCompressValue(mustCompressValue MustCompressValue) {
	r.mustCompressValue = mustCompressValue
}

func (r MustCompressRead) getMustCompressValue() (mustCompressValue MustCompressValue) {
	mustCompressValue = r.mustCompressValue
	return
}

func (r MustCompressRead) genMustCompressValue() (mustCompressValue MustCompressValue, err error) {
	mustCompressStr := r.getMustCompressString()

	val, errVal := ConvertStringToMustCompressValue(mustCompressStr)
	if errVal != nil {
		err = errVal
		return
	}

	mustCompressValue = val
	return
}

func (r *MustCompressRead) setMustCompressString(mustCompressString string) {
	r.mustCompressString = mustCompressString
}

func (r MustCompressRead) getMustCompressString() (mustCompressString string) {
	mustCompressString = r.mustCompressString
	return
}

func (r MustCompressRead) getMustCompressStringFromRedis() (mustCompressString string, err error) {
	builder := r.getBuilder()
	redisGet := builder.RedisGet

	key := redisMustCompressKey

	resGet, errGet := redisGet.Get(context.Background(), key).Result()
	if errGet != nil {
		err = errGet
		return
	}

	mustCompressString = resGet
	return
}

func (bac *MustCompressRead) setErr(methodName string, err error) {
	formatErr := fmt.Errorf("redismustcompress.NewMustCompressRead.%s: %s", methodName, err)
	bac.err = formatErr
}

func (mcs *MustCompressRead) clearErr() {
	mcs.err = nil
}

func (mcs MustCompressRead) getErr() (err error) {
	err = mcs.err
	return
}
