package redismustcompress

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type IReader interface {
	Read()
	IsFound() (isFound bool)
	MustCompress() (mustCompress bool)
	Err() (err error)
}

type ISaver interface {
	Save(mustCompress bool)
	Err() (err error)
}

type BuilderCompareMustCompress struct {
	Reader IReader
	Saver  ISaver
}

type MustCompressCompare struct {
	builder         BuilderCompareMustCompress
	envValue        bool
	redisValue      bool
	isRedisKeyFound bool
	mustCompress    bool
	actionSave      bool
	err             error
}

func NewCompareMustCompress(builder BuilderCompareMustCompress) (mcc *MustCompressCompare) {
	mcc = new(MustCompressCompare)
	// builder
	mcc.setBuilder(builder)
	return
}

func (mcc *MustCompressCompare) Compare() {
	// clearErr first
	mcc.clearErr()

	// envValue
	envValue := mcc.getEnvValueFromViper()
	mcc.setEnvValue(envValue)

	// redisValue, isRedisKeyFound
	redisValue, isRedisKeyFound, errRedisValue := mcc.getRedisValueFromRedis()
	if errRedisValue != nil {
		mcc.setErr("getRedisValueFromRedis()", errRedisValue)
		return
	}
	mcc.setIsRedisKeyFound(isRedisKeyFound)
	mcc.setRedisValue(redisValue)

	// compare, mustCompress, actionSave
	actionSave, mustCompress, errCompare := mcc.compareEnvAndRedisValue()
	if errCompare != nil {
		mcc.setErr("compareEnvAndRedisValue()", errCompare)
		return
	}
	mcc.setMustCompress(mustCompress)
	mcc.setActionSave(actionSave)

	// actionSave
	mcc.runActionSave()
}

func (mcc MustCompressCompare) Err() (err error) {
	err = mcc.getErr()
	return
}

func (mcc MustCompressCompare) MustCompress() (mustCompress bool) {
	mustCompress = mcc.getMustCompress()
	return
}

func (mcc *MustCompressCompare) setBuilder(builder BuilderCompareMustCompress) {
	mcc.builder = builder
}

func (mcc MustCompressCompare) getBuilder() (builder BuilderCompareMustCompress) {
	builder = mcc.builder
	return
}

func (mcc *MustCompressCompare) setEnvValue(envValue bool) {
	mcc.envValue = envValue
}

func (mcc MustCompressCompare) getEnvValue() (envValue bool) {
	envValue = mcc.envValue
	return
}

func (mcc MustCompressCompare) getEnvValueFromViper() (envValue bool) {
	envValue = viper.GetBool(viperRedisCompressionKey)
	return
}

func (mcc *MustCompressCompare) setRedisValue(redisValue bool) {
	mcc.redisValue = redisValue
}

func (mcc MustCompressCompare) getRedisValue() (redisValue bool) {
	redisValue = mcc.redisValue
	return
}

func (mcc MustCompressCompare) getRedisValueFromRedis() (redisValue bool, isFound bool, err error) {
	builder := mcc.builder
	reader := builder.Reader

	reader.Read()
	if reader.Err() != nil {
		err = reader.Err()
		return
	}

	redisValue = reader.MustCompress()
	isFound = reader.IsFound()
	return
}

func (mcc *MustCompressCompare) setIsRedisKeyFound(isRedisKeyFound bool) {
	mcc.isRedisKeyFound = isRedisKeyFound
}

func (mcc MustCompressCompare) getIsRedisKeyFound() (isRedisKeyFound bool) {
	isRedisKeyFound = mcc.isRedisKeyFound
	return
}

func (mcc *MustCompressCompare) setMustCompress(mustCompress bool) {
	mcc.mustCompress = mustCompress
}

func (mcc MustCompressCompare) getMustCompress() (mustCompress bool) {
	mustCompress = mcc.mustCompress
	return
}

/*
compareEnvAndRedisValue

Compare mustCompress from redis and from env
*/
func (mcs MustCompressCompare) compareEnvAndRedisValue() (actionSave bool, mustCompress bool, err error) {
	envValue := mcs.getEnvValue()
	isRedisKeyFound := mcs.getIsRedisKeyFound()
	redisValue := mcs.getRedisValue()

	// mustCompress from redis is not found
	if !isRedisKeyFound {
		actionSave = true
		mustCompress = envValue
		return
	}

	// mustCompress from redis == mustCompress from env
	if redisValue == envValue {
		mustCompress = redisValue
		return
	}

	// mustCompress from redis != mustCompress from env
	envValueStr := strconv.FormatBool(envValue)
	redisValueStr := strconv.FormatBool(redisValue)

	err = fmt.Errorf("%[1]s environment has been changed to %[2]s (,but state recorded in redis was %[3]s). please flush all your cache before change the environment. you can also undo the changes", viperRedisCompressionKey, envValueStr, redisValueStr)
	return
}

func (mcc *MustCompressCompare) setActionSave(actionSave bool) {
	mcc.actionSave = actionSave
}

func (mcc MustCompressCompare) getActionSave() (actionSave bool) {
	actionSave = mcc.actionSave
	return
}

func (mcc MustCompressCompare) runActionSave() (err error) {
	actionSave := mcc.getActionSave()
	if !actionSave {
		return
	}

	builder := mcc.getBuilder()
	saver := builder.Saver

	mustCompress := mcc.getMustCompress()

	saver.Save(mustCompress)
	if saver.Err() != nil {
		err = saver.Err()
		return
	}

	return
}

func (bac *MustCompressCompare) setErr(methodName string, err error) {
	formatErr := fmt.Errorf("redismustcompress.NewMustCompressCompare.%s: %s", methodName, err)
	bac.err = formatErr
}

func (mcs *MustCompressCompare) clearErr() {
	mcs.err = nil
}

func (mcs MustCompressCompare) getErr() (err error) {
	err = mcs.err
	return
}

// Summary:
// ==========
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
