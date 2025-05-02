package redismustcompress

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisSet interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type BuilderMustCompressSave struct {
	RedisSet IRedisSet
}

type MustCompressSave struct {
	builder            BuilderMustCompressSave
	mustCompress       bool
	mustCompressValue  MustCompressValue
	mustCompressString string
	err                error
}

func NewMustCompressSave(builder BuilderMustCompressSave) (mcs *MustCompressSave) {
	mcs = new(MustCompressSave)
	// builder
	mcs.setBuilder(builder)
	return
}

func (mcs *MustCompressSave) Save(mustCompress bool) {
	// clearErr first
	mcs.clearErr()

	// mustCompress
	mcs.setMustCompress(mustCompress)
	defer mcs.clearMustCompress()

	// mustCompressValue
	mustCompressVal := mcs.genMustCompressValue()
	mcs.setMustCompressValue(mustCompressVal)

	// mustCompressString
	mustCompressString := mcs.genMustCompressString()
	mcs.setMustCompressString(mustCompressString)

	// save mustCompressString
	errSave := mcs.saveMustCompressStringToRedis()
	mcs.setErr("saveMustCompressStringToRedis()", errSave)
}

func (mcs MustCompressSave) Me() MustCompressSave {
	return mcs
}

func (mcs MustCompressSave) Err() (err error) {
	err = mcs.getErr()
	return
}

func (r *MustCompressSave) setBuilder(builder BuilderMustCompressSave) {
	r.builder = builder
}

func (r MustCompressSave) getBuilder() (builder BuilderMustCompressSave) {
	builder = r.builder
	return
}

func (r *MustCompressSave) setMustCompress(mustCompress bool) {
	r.mustCompress = mustCompress
}

func (r *MustCompressSave) clearMustCompress() {
	r.mustCompress = false
}

func (r MustCompressSave) getMustCompress() (mustCompress bool) {
	mustCompress = r.mustCompress
	return
}

func (r *MustCompressSave) setMustCompressValue(mustCompressValue MustCompressValue) {
	r.mustCompressValue = mustCompressValue
}

func (r MustCompressSave) getMustCompressValue() (mustCompressValue MustCompressValue) {
	mustCompressValue = r.mustCompressValue
	return
}

func (r MustCompressSave) genMustCompressValue() (mustCompressValue MustCompressValue) {
	mustCompress := r.getMustCompress()
	mustCompressValue = NewMustCompressValue(mustCompress)
	return
}

func (r *MustCompressSave) setMustCompressString(mustCompressString string) {
	r.mustCompressString = mustCompressString
}

func (r MustCompressSave) getMustCompressString() (mustCompressString string) {
	mustCompressString = r.mustCompressString
	return
}

func (r MustCompressSave) genMustCompressString() (mustCompressString string) {
	mustCompressValue := r.getMustCompressValue()
	mustCompressString = mustCompressValue.String()
	return
}

func (r MustCompressSave) saveMustCompressStringToRedis() (err error) {
	builder := r.getBuilder()
	redisSet := builder.RedisSet

	mustCompressString := r.getMustCompressString()

	key := redisMustCompressKey

	errSet := redisSet.Set(context.Background(), key, mustCompressString, 0).Err()
	if errSet != nil {
		err = errSet
		return
	}

	return
}

func (bac *MustCompressSave) setErr(methodName string, err error) {
	formatErr := fmt.Errorf("redismustcompress.NewMustCompressSave.%s: %s", methodName, err)
	bac.err = formatErr
}

func (mcs *MustCompressSave) clearErr() {
	mcs.err = nil
}

func (mcs MustCompressSave) getErr() (err error) {
	err = mcs.err
	return
}
