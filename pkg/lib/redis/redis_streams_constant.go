package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/terra-discover/bbcrs-helper-lib/pkg/lib"
)

func GetConsumerName() (consumerName string, err error) {
	envVar := "APP_NAME"

	appName := viper.GetString(envVar)
	if lib.IsEmptyStr(appName) {
		err = fmt.Errorf("services.ConsumerName(): cannot specify consumer name. Env %s is empty or undefined", envVar)
		return
	}

	consumerName = appName
	return
}

const (
	BookingGroupName string = "booking_group"
)

const (
	TrackBookingStreamName string = "track:booking"
)

/*
TransportType

Used in RedisStreamTransport
*/
const (
	BookingNotifiedTransportType string = "booking_notified"
)

/*
BookingNotifiedTransport

TransportType = BookingNotifiedTransportType
*/
type BookingNotifiedTransport struct {
	AgentID     uuid.UUID `json:"agent_id"`
	CorporateID uuid.UUID `json:"corporate_id"`
	ProposalID  uuid.UUID `json:"proposal_id"`
	// SendTo      ListReceiver `json:"send_to"`
	// ReceivedBy  ListReceiver `json:"received_by"` // will be filled by receiver
	// ValidFrom   strfmt.DateTime `json:"valid_from"`
	// ValidTo     strfmt.DateTime `json:"valid_to"`
}

func (tbm BookingNotifiedTransport) JsonString() (jsonString string, err error) {
	bte, errMarshal := lib.JSONMarshal(tbm)
	if errMarshal != nil {
		err = errMarshal
		return
	}

	jsonString = string(bte)
	return
}

/*
redis streams special IDs

Source: [Special IDs in the streams API](https://redis.io/docs/data-types/streams/#special-ids-in-the-streams-api)
*/
const (
	DollarSign string = "$"
	ZeroNumber string = "0"
	GtSign     string = ">"
	StarSign   string = "*"
	MinusSign  string = "-"
	PlusSign   string = "+"
)

type CompressTool string

const (
	NoneCompressTool CompressTool = "none"
	GzipCompressTool CompressTool = "gzip"
)

func (cmt CompressTool) String() string {
	return string(cmt)
}

func GenCompressTool(input string) (compressTool CompressTool, err error) {
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case NoneCompressTool.String():
		{
			compressTool = NoneCompressTool
			break
		}
	case GzipCompressTool.String():
		{
			compressTool = GzipCompressTool
			break
		}
	default:
		{
			err = fmt.Errorf("services.GenCompressTool: compress tool %s is unknown", input)
			break
		}
	}

	return
}

type RedisStreamTransport struct {
	TransportType string `json:"transport_type"`
	CompressTool  string `json:"compress_tool"`
	Data          string `json:"data"`
}

func (rst RedisStreamTransport) JsonString() (result string, err error) {
	bte, errMarshal := lib.JSONMarshal(rst)
	if errMarshal != nil {
		err = errMarshal
		return
	}

	result = string(bte)
	return
}

func (rst RedisStreamTransport) MapInterface() (result map[string]interface{}, err error) {
	err = lib.Merge(rst, &result)
	return
}
