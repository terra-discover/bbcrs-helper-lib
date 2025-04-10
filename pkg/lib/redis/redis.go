package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/terra-discover/bbcrs-helper-lib/pkg/lib"
	"github.com/terra-discover/bbcrs-helper-lib/pkg/lib/redis/redismustcompress"
)

// REDIS null if not initialized
var REDIS *redis.Client

// InitRedis initialize redis connection
func InitRedis() {
	if nil == REDIS {
		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		if redisHost != "" {
			REDIS = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
				Password: viper.GetString("REDIS_PASS"),
				DB:       viper.GetInt("REDIS_INDEX"),
			})
		}
	}
}

// SetCachingRedis func
func SetCachingRedis(rdb *redis.Client, datas map[string]map[string]interface{}) {
	repo := NewRedisRepository(rdb)
	re := regexp.MustCompile(`[0-9]+$`)
	for k, v := range datas {
		key := re.ReplaceAll([]byte(k), []byte(``))
		err := repo.Set(string(key), lib.ConvertJsonToStr(v["values"]), 0)
		if err != nil {
			fmt.Printf("unable to SET data. error: %v", err)
		}
	}
}

// RedisRepository represent the repositories
type RedisRepository interface {
	Set(key string, value string, exp time.Duration) error
	MSet(mapKeyValues map[string]string, exp time.Duration) error
	Get(key string) (string, error)
	MGet(keys []string) (map[string]string, error)
	IsExist(keys ...string) (bool, error)
	Del(key string) (int, error)
	GetDel(key string) (string, error)
	AppendStartList(key string, value ...string) (int64, error)
	AppendEndList(key string, value ...string) (int64, error)
	GetList(key string, start int64, end int64) ([]string, error)
	RemoveMatchFromList(key string, count int64, matchValue string) (int64, error)
	LeftPopCountList(key string, count uint) ([]string, error)
	compress(key, val string) (result string, err error)
	decompress(val string) (result string, err error)
}

// repository represent the repository model
type redisRepository struct {
	Client       *redis.Client
	reader       redismustcompress.MustCompressRead
	saver        redismustcompress.MustCompressSave
	mustCompress bool
	err          error
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRedisRepository(client *redis.Client) (r *redisRepository) {
	if r == nil {
		r = new(redisRepository)
	}

	// Client
	r.setClient(client)
	return
}

func (r redisRepository) MustCompress() (mustCompress bool) {
	mustCompress = r.getMustCompress()
	return
}

func (r redisRepository) Err() (err error) {
	err = r.getErr()
	return
}

func (r *redisRepository) NewSession() (newR *redisRepository) {
	newR = r

	// clearErr first
	newR.clearErr()

	// reader
	reader := newR.genReader()
	newR.setReader(reader)

	// saver
	saver := newR.genSaver()
	newR.setSaver(saver)

	// mustCompress
	mustCompress, errMustCompress := newR.genMustCompress()
	if errMustCompress != nil {
		newR.setErr(errMustCompress)
		return
	}
	newR.setMustCompress(mustCompress)

	return
}

func (r *redisRepository) setClient(client *redis.Client) {
	r.Client = client
}

func (r redisRepository) getClient() (client *redis.Client) {
	client = r.Client
	return
}

func (r *redisRepository) setReader(reader redismustcompress.MustCompressRead) {
	r.reader = reader
}

func (r redisRepository) getReader() (reader redismustcompress.MustCompressRead) {
	reader = r.reader
	return
}

func (r redisRepository) genReader() (reader redismustcompress.MustCompressRead) {
	client := r.getClient()

	builderReader := redismustcompress.BuilderMustCompressRead{
		RedisGet: client,
	}
	reader = redismustcompress.NewMustCompressRead(builderReader).Me()
	return
}

func (r *redisRepository) setSaver(saver redismustcompress.MustCompressSave) {
	r.saver = saver
}

func (r redisRepository) getSaver() (saver redismustcompress.MustCompressSave) {
	saver = r.saver
	return
}

func (r redisRepository) genSaver() (saver redismustcompress.MustCompressSave) {
	client := r.getClient()

	builderSaver := redismustcompress.BuilderMustCompressSave{
		RedisSet: client,
	}
	saver = redismustcompress.NewMustCompressSave(builderSaver).Me()
	return
}

func (r *redisRepository) setMustCompress(mustCompress bool) {
	r.mustCompress = mustCompress
}

func (r redisRepository) getMustCompress() (mustCompress bool) {
	mustCompress = r.mustCompress
	return
}

func (r redisRepository) genMustCompress() (mustCompress bool, err error) {
	reader := r.getReader()
	saver := r.getSaver()

	builderCompare := redismustcompress.BuilderCompareMustCompress{
		Reader: &reader,
		Saver:  &saver,
	}
	cmp := redismustcompress.NewCompareMustCompress(builderCompare)

	cmp.Compare()
	if cmp.Err() != nil {
		err = cmp.Err()
		return
	}

	mustCompress = cmp.MustCompress()
	return
}

func (mcs *redisRepository) setErr(err error) {
	mcs.err = err
}

func (mcs *redisRepository) clearErr() {
	mcs.err = nil
}

func (mcs redisRepository) getErr() (err error) {
	err = mcs.err
	return
}

// Set attaches the redis repository and set the data
func (r *redisRepository) Set(key, value string, exp time.Duration) error {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return r.Err()
	}

	valCompress, errCompress := r.compress(key, value)
	if errCompress != nil {
		return errCompress
	}

	return r.Client.Set(context.Background(), key, valCompress, exp).Err()
}

/*
MSet

MSet params map[string]string => map[key]value

-Using pipeline for including expire time

Source: https://redis.uptrace.dev/guide/go-redis-pipelines.html
*/
func (r *redisRepository) MSet(mapKeyValues map[string]string, exp time.Duration) error {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return r.Err()
	}

	// Need pipeline for including expire time
	cmds, err := r.Client.Pipelined(context.TODO(), func(rd redis.Pipeliner) error {
		for key := range mapKeyValues {
			val := mapKeyValues[key]

			valCompress, errCompress := r.compress(key, val)
			if errCompress != nil {
				return errCompress
			}

			rd.Set(context.TODO(), key, valCompress, exp)
		}

		return nil
	})

	if err != nil {
		return err
	}

	arrErr := []string{}

	for idxCmd := range cmds {
		itemCmd := cmds[idxCmd]

		// Casting type as result of redis.Set() method
		res, ok := itemCmd.(*redis.StatusCmd)
		if !ok {
			continue
		}

		if res.Err() != nil {
			loc := fmt.Sprintf("index: %d", idxCmd)
			mess := fmt.Sprintf("%s, message: %s", loc, res.Err().Error())
			arrErr = append(arrErr, mess)
			continue
		}
	}

	if len(arrErr) > 0 {
		mess := strings.Join(arrErr, " ; /n")
		return errors.New(mess)
	}

	return nil
}

// Get attaches the redis repository and get the data
func (r *redisRepository) Get(key string) (string, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return "", r.Err()
	}

	get := r.Client.Get(context.Background(), key)

	val, errGet := get.Result()
	if errGet != nil {
		return "", errGet
	}

	return r.decompress(val)
}

/*
MGet

MGet with returning map[key]value => map[string]string

If value of key is empty or key not found, will return map[key]empty_string

-Keys order and return values order are always guaranteed same

Source: https://github.com/redis/redis/issues/4647#issuecomment-362502460
*/
func (r *redisRepository) MGet(keys []string) (map[string]string, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return nil, r.Err()
	}

	finalRes := make(map[string]string)

	res := r.Client.MGet(context.TODO(), keys...)
	if res.Err() != nil {
		return nil, res.Err()
	}

	values := res.Val()

	arrErr := []string{}

	for idxVal := range values {
		itemVal := values[idxVal]
		itemKey := keys[idxVal]

		strVal, ok := itemVal.(string)
		if !ok {
			finalRes[itemKey] = ""
		} else {
			valDecompress, errDecompress := r.decompress(strVal)
			if errDecompress != nil {
				loc := fmt.Sprintf("index: %s", itemKey)
				mess := fmt.Sprintf("%s, message: %s", loc, errDecompress.Error())
				arrErr = append(arrErr, mess)
				continue
			}

			finalRes[itemKey] = valDecompress
		}
	}

	if len(arrErr) > 0 {
		mess := strings.Join(arrErr, " ; /n")
		return nil, errors.New(mess)
	}

	return finalRes, nil
}

func (r *redisRepository) IsExist(keys ...string) (bool, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return false, r.Err()
	}

	result := r.Client.Exists(context.Background(), keys...)
	return result.Val() > 0, result.Err()
}

// Delete the data by key
func (r *redisRepository) Del(key string) (count int, err error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return 0, r.Err()
	}

	res := r.Client.Del(context.Background(), key)
	if res.Err() != nil {
		err = res.Err()
		return
	}

	count = int(res.Val())
	return
}

// Delete the data by key and return values
func (r *redisRepository) GetDel(key string) (string, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return "", r.Err()
	}

	get := r.Client.GetDel(context.Background(), key)
	if get.Err() != nil {
		return "", get.Err()
	}

	val := get.Val()

	return r.decompress(val)
}

// Append from start list
func (r *redisRepository) AppendStartList(key string, values ...string) (int64, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return 0, r.Err()
	}

	newValues := []string{}

	// Remove first if exist
	for idxVal := range values {
		itemVal := values[idxVal]

		valCompress, errCompress := r.compress(key, itemVal)
		if errCompress != nil {
			return 0, errCompress
		}

		// append
		newValues = append(newValues, valCompress)
	}

	// Append
	get := r.Client.LPush(context.Background(), key, newValues)
	return get.Result()
}

// Append in end list
func (r *redisRepository) AppendEndList(key string, values ...string) (int64, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return 0, r.Err()
	}

	newValues := []string{}

	// Remove first if exist
	for idxVal := range values {
		itemVal := values[idxVal]

		valCompress, errCompress := r.compress(key, itemVal)
		if errCompress != nil {
			return 0, errCompress
		}

		// append
		newValues = append(newValues, valCompress)
	}

	// Append
	get := r.Client.RPush(context.Background(), key, newValues)
	return get.Result()
}

// params from = 0 and end = -1 , will return all entire list
func (r *redisRepository) GetList(key string, start int64, end int64) ([]string, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return []string{}, r.Err()
	}

	get := r.Client.LRange(context.Background(), key, start, end)
	if get.Err() != nil {
		return []string{}, get.Err()
	}

	values := get.Val()

	result := []string{}

	for idxVal := range values {
		itemVal := values[idxVal]

		valDecompress, errDecompress := r.decompress(itemVal)
		if errDecompress != nil {
			return []string{}, errDecompress
		}

		result = append(result, valDecompress)
	}

	return result, nil
}

// Removes the first count occurrences of elements equal to element from the list stored at key. The count argument influences the operation in the following ways:
// count > 0: Remove elements equal to element moving from head to tail.
// count < 0: Remove elements equal to element moving from tail to head.
// count = 0: Remove all elements equal to element.
// For example, LREM list -2 "hello" will remove the last two occurrences of "hello" in the list stored at list.
// Note that non-existing keys are treated like empty lists, so when key does not exist, the command will always return 0.
// Source: https://redis.io/commands/lrem/
func (r *redisRepository) RemoveMatchFromList(key string, count int64, matchValue string) (int64, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return 0, r.Err()
	}

	valCompress, errCompress := r.compress(key, matchValue)
	if errCompress != nil {
		return 0, errCompress
	}

	remove := r.Client.LRem(context.Background(), key, count, valCompress)
	return remove.Result()
}

/*
LeftPopCountList

-Remove list from left with count

Example:
  - list: ["a", "b", "c", "d"];
  - LeftPopCountList with count = 2;
  - Result: ["c", "d"];
*/
func (r *redisRepository) LeftPopCountList(key string, count uint) ([]string, error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return []string{}, r.Err()
	}

	leftPop := r.Client.LPopCount(context.Background(), key, int(count))
	return leftPop.Result()
}

func (r redisRepository) compress(key, val string) (result string, err error) {
	setErr := func(err error) (resultErr error) {
		if err == nil {
			return
		}

		oldErr := err.Error()
		resultErr = fmt.Errorf("redis compress: %s", oldErr)
		return
	}

	mustCompress := r.mustCompress
	if !mustCompress {
		result = val
		return
	}

	header := lib.GzipHeader{
		Name:    key,
		ModTime: time.Now().UTC(),
	}

	resCompress, errCompress := lib.CompressGzipString(val, header)
	if errCompress != nil {
		err = setErr(errCompress)
		return
	}

	result = resCompress
	return
}

func (r redisRepository) decompress(val string) (result string, err error) {
	setErr := func(err error) (resultErr error) {
		if err == nil {
			return
		}

		oldErr := err.Error()
		resultErr = fmt.Errorf("redis decompress: %s. If you met any decompressing issue, please make sure first that you are using compress on set key value, in ex: check env REDIS_COMPRESSION=true, or check compressing tools of set key value must same with decompressing, else make sure your operating system is supported for the compressing tools", oldErr)
		return
	}

	mustCompress := r.mustCompress
	if !mustCompress {
		result = val
		return
	}

	resDecompress, _, errDecompress := lib.DecompressGzipString(val)
	if errDecompress != nil {
		err = setErr(errDecompress)
		return
	}

	result = resDecompress
	return
}
