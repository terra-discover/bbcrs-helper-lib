package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/terra-discover/bbcrs-helper-lib/pkg/lib"
)

// XGroupCreate attaches the redis repository and set the data
func (r *redisRepository) XGroupCreate(stream, group string) (result string, err error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return "", r.Err()
	}

	resXGroup, errXGroup := r.Client.XGroupCreateMkStream(context.Background(), stream, group, DollarSign).Result()
	if errXGroup != nil {
		err = fmt.Errorf("services.XGroupCreate(): %s", errXGroup)
		return
	}

	result = resXGroup
	return
}

func (r *redisRepository) XAdd(stream, transportType, value string) (result string, err error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return "", r.Err()
	}

	compressTool := NoneCompressTool

	mustCompress := r.MustCompress()
	if mustCompress {
		compressTool = GzipCompressTool
	}

	resCompress, errCompress := r.compress(transportType, value)
	if errCompress != nil {
		err = fmt.Errorf("services.XAdd(): %s", errCompress)
		return
	}

	trans := RedisStreamTransport{
		TransportType: transportType,
		CompressTool:  compressTool.String(),
		Data:          resCompress,
	}

	mapValues, errMapValues := trans.MapInterface()
	if errMapValues != nil {
		err = fmt.Errorf("services.XAdd(): %s", errMapValues)
		return
	}

	resXAdd, errXAdd := r.Client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: mapValues,
	}).Result()
	if errXAdd != nil {
		err = fmt.Errorf("services.XAdd(): %s", errXAdd)
		return
	}

	result = resXAdd
	return
}

/*
XReadGroup

result: map[stream id]value

1 stream must only have 1 transport type to minimalize complexity

params: mapStreamNameID (map[string]string) => map[stream name]stream id

Stream id has 2 kind:
  - specific stream id
  - special id

The id can be special IDs, such as:
  - zeroNumber "0": return pending list stream from Pending List (PEL)
  - gt sign ">": return actual list stream which unread and not in Pending List (PEL)

Not all special IDs can be used in XReadGroup() function.

Source:
  - [Special IDs in the streams API](https://redis.io/docs/data-types/streams/#special-ids-in-the-streams-api)
*/
func (r *redisRepository) XReadGroup(mapStreamNameID map[string]string, group string) (result map[string]string, transportType string, err error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return nil, "", r.Err()
	}

	consumerName, errConsumer := GetConsumerName()
	if errConsumer != nil {
		err = fmt.Errorf("services.XReadGroup().GetConsumerName(): %s", errConsumer)
		return
	}

	streams := prepareStreams(mapStreamNameID)

	resXReadGroup, errXReadGroup := r.Client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumerName,
		Streams:  streams,
		Block:    1 * time.Second,
	}).Result()
	if errXReadGroup != nil {
		err = fmt.Errorf("services.XReadGroup().XReadGroup(): %s", errXReadGroup)
		return
	}

	result = make(map[string]string)

	if len(resXReadGroup) == 0 {
		return
	}

	resMessages := resXReadGroup[0].Messages
	if len(resMessages) == 0 {
		return
	}

	for idxMess := range resMessages {
		itemMess := resMessages[idxMess]

		id := itemMess.ID
		mapValue := itemMess.Values

		transType, value, errValue := r.xreadMapValue(mapValue)
		if errValue != nil {
			err = fmt.Errorf("services.XReadGroup().xreadMapValue(): resMessages id %s: %s", id, errValue)
			break
		}

		result[id] = value

		if lib.IsEmptyStr(transportType) {
			transportType = transType
		}
	}

	if err != nil {
		return
	}

	return
}

func prepareStreams(mapStreamNameID map[string]string) (result []string) {
	listStreamName := []string{}
	listStreamID := []string{}

	for streamName := range mapStreamNameID {
		streamID := mapStreamNameID[streamName]

		listStreamName = append(listStreamName, streamName)
		listStreamID = append(listStreamID, streamID)
	}

	// Append result must be ordered
	// Add list stream name first, then append list stream id
	result = listStreamName
	result = append(result, listStreamID...)
	return
}

func (r *redisRepository) xreadMapValue(mapValue map[string]interface{}) (transportType, value string, err error) {
	if mapValue == nil {
		return
	}

	// merge to RedisStreamTransport struct
	streamTransport := RedisStreamTransport{}

	errMerge := lib.Merge(mapValue, &streamTransport)
	if errMerge != nil {
		err = fmt.Errorf("services.xreadMapValue().lib.Merge(): %s", errMerge)
		return
	}

	rawTransportType := streamTransport.TransportType
	rawCompressTool := streamTransport.CompressTool
	rawData := streamTransport.Data

	// get compress tool
	compressTool, errCompressTool := GenCompressTool(rawCompressTool)
	if errCompressTool != nil {
		err = fmt.Errorf("services.xreadMapValue().GenCompressTool(): %s", errCompressTool)
		return
	}

	// decompress
	if compressTool == GzipCompressTool {
		resDecompress, errDecompress := r.decompress(rawData)
		if errDecompress != nil {
			err = fmt.Errorf("services.xreadMapValue().decompress(): %s", errDecompress)
			return
		}

		rawData = resDecompress
	}

	transportType = rawTransportType
	value = rawData
	return
}

/*
XInfoGroup
*/
func (r *redisRepository) XInfoGroups(stream string) (result []redis.XInfoGroup, isFound bool, err error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return []redis.XInfoGroup{}, false, r.Err()
	}

	const notFoundErrSubstr = "no such key"

	resXInfoGroups, errXInfoGroups := r.Client.XInfoGroups(context.Background(), stream).Result()
	if errXInfoGroups != nil {
		mess := errXInfoGroups.Error()
		if strings.Contains(mess, notFoundErrSubstr) {
			return
		}

		err = fmt.Errorf("services.XInfoGroup(): %s", errXInfoGroups)
		return
	}

	isFound = true
	result = resXInfoGroups
	return
}

/*
XAck

Acknowledge stream
*/
func (r *redisRepository) XAck(stream, group string, streamIDs []string) (result int64, err error) {
	// start session
	r.NewSession()
	if r.Err() != nil {
		return 0, r.Err()
	}

	resXAck, errXAck := r.Client.XAck(context.Background(), stream, group, streamIDs...).Result()
	if errXAck != nil {
		err = fmt.Errorf("services.XAck(): %s", errXAck)
		return
	}

	result = resXAck
	return
}
