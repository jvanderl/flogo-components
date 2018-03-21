package encoder

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-jl-encoder")

const (
	actionEncode  = "ENCODE"
	actionDecode  = "DECODE"
	encoderBase64 = "BASE64"
	encoderBase32 = "BASE32"
	encoderHex    = "HEX"
	ivEncoder     = "encoder"
	ivAction      = "action"
	ivInput       = "input"
	ovResult      = "result"
	ovStatus      = "status"
)

// encoderActivity is an Activity that is used to encode and decode strings
// inputs : {type,action,input}
// outputs: {result,status}
type encoderActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new encoderActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &encoderActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *encoderActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *encoderActivity) Eval(context activity.Context) (done bool, err error) {

	encoder := strings.ToUpper(context.GetInput(ivEncoder).(string))
	action := strings.ToUpper(context.GetInput(ivAction).(string))
	input := context.GetInput(ivInput).(string)
	log.Debugf("input is: %v", input)
	status := ""
	result := ""
	var errEncode error
	switch action {
	case actionEncode:
		log.Debug("Starting Encode")
		switch encoder {
		case encoderBase32:
			log.Debug("Using Base32 encoder")
			result, errEncode = encodeBase32(input)
		case encoderBase64:
			log.Debug("Using Base64 encoder")
			result, errEncode = encodeBase64(input)
		case encoderHex:
			log.Debug("Using Hex encoder")
			result, errEncode = encodeHex(input)
		default:
			log.Debugf("Invalid encoder: %v", encoder)
		}
	case actionDecode:
		log.Debug("Starting Decode")
		switch encoder {
		case encoderBase32:
			log.Debug("Using Base32 decoder")
			result, errEncode = decodeBase32(input)
		case encoderBase64:
			log.Debug("Using Base64 decoder")
			result, errEncode = decodeBase64(input)
		case encoderHex:
			log.Debug("Using Hex decoder")
			result, errEncode = decodeHex(input)
		default:
			log.Debugf("Invalid encoder: %v", encoder)
		}
	default:
		log.Debugf("Invalid Action: %v", action)
	}

	if errEncode != nil {
		status = errEncode.Error()
	} else {
		status = "OK"
	}
	context.SetOutput(ovResult, result)
	context.SetOutput(ovStatus, status)
	return true, nil
}

func encodeBase32(input string) (string, error) {
	output := base32.StdEncoding.EncodeToString([]byte(input))
	return output, nil
}

func decodeBase32(input string) (string, error) {
	data, err := base32.StdEncoding.DecodeString(input)
	return string(data), err
}

func encodeBase64(input string) (string, error) {
	output := base64.StdEncoding.EncodeToString([]byte(input))
	return output, nil
}

func decodeBase64(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	return string(data), err
}

func encodeHex(input string) (string, error) {
	output := hex.EncodeToString([]byte(input))
	return output, nil
}

func decodeHex(input string) (string, error) {
	data, err := hex.DecodeString(input)
	return string(data), err
}
