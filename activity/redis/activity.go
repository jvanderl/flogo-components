package redis

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/go-redis/redis"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-redis")

const (
	server		= "server"
	password	= "password"
	database	= "database"
	operation	= "operation"
	key				= "key"
	value			= "value"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	serverInput := context.GetInput(server)

	ivServer, ok := serverInput.(string)
	if !ok {
		context.SetOutput("result", "SERVER_NOT_SET")
		return true, fmt.Errorf("server not set")
	}



	valueInput := context.GetInput(value)

	ivValue, ok := valueInput.(string)
	if !ok {
		ivValue = ""
	}

	operationInput := context.GetInput(operation)

	ivOperation, ok := operationInput.(string)
	if !ok {
		context.SetOutput("result", "OPERATION_NOT_SET")
		return true, fmt.Errorf("operation not set")
	}

	keyInput := context.GetInput(key)

	ivKey, ok := keyInput.(string)
	if !ok && ivOperation != "PING" {
		context.SetOutput("result", "KEY_NOT_SET")
		return true, fmt.Errorf("key not set")
	}

	ivDatabase, ok := context.GetInput(database).(int)
	if !ok {
		//User not set, use default
		ivDatabase = 0
	}

	passwordInput := context.GetInput(password)

	ivPassword, ok := passwordInput.(string)
	if !ok {
		//Password not set, use default
		ivPassword = ""
	}

	log.Debugf("Connecting to Redis server")
	client := redis.NewClient(&redis.Options{
		Addr:     ivServer,
		Password: ivPassword,
		DB:       ivDatabase,
	})

	switch ivOperation {
	case "GET":
		val, err := client.Get(ivKey).Result()
		if err == redis.Nil {
			context.SetOutput("result", "KEY_NOT_FOUND")
			log.Debugf("Key was not found: %v", err)
			return true, nil
		} else if err != nil {
			context.SetOutput("result", "ERROR_GET_VALUE")
			return true, err
		}
		context.SetOutput("result", val)
	case "KEYS":
		val, err := client.Keys(ivKey).Result()
		if err == redis.Nil {
			context.SetOutput("result", "KEY_NOT_FOUND")
			log.Debugf("Keys not found: %v", err)
			return true, nil
		} else if err != nil {
			context.SetOutput("result", "ERROR_GET_KEYS")
			return true, err
		}
		context.SetOutput("result", val)
/*	case "MGET":
		val, err := client.Mget(ivKey).Result()
		if err == redis.Nil {
			context.SetOutput("result", "KEY_NOT_FOUND")
			log.Debugf("Key was not found: %v", err)
			return true, nil
		} else if err != nil {
			context.SetOutput("result", "ERROR_GET_VALUE")
			return true, err
		}
		context.SetOutput("result", val) */
	case "SET":
		err := client.Set(ivKey, ivValue, 0).Err()
		if err != nil {
			context.SetOutput("result", "ERROR_SET_VALUE")
			return true, err
		}
		context.SetOutput("result", "OK")
/*	case "MSET":
		pairs := []interface{}
		pairs[0]="unos"
		pairs[1]="dos"
		var status string = ""
		err := client.MSet(pairs, status).Err()
		if err != nil {
			context.SetOutput("result", "ERROR_SET_MULTI")
			log.Errorf ("Error: %v", err)
			return true, err
		}
		context.SetOutput("result", "OK") */
	case "DEL":
		err := client.Del(ivKey).Err()
		if err != nil {
			context.SetOutput("result", "ERROR_DEL_KEY")
			return true, err
		}
		context.SetOutput("result", "OK")
	case "PING":
		pong, err := client.Ping().Result()
		if err != nil {
			context.SetOutput("result", "ERROR_PING")
			return true, err
		}
		context.SetOutput("result", pong)
	default:
		context.SetOutput("result", "UNKNOWN_OPERATOR")
	}

	client.Close()

	log.Debugf("Redis server disconnected")
	return true, nil
}
