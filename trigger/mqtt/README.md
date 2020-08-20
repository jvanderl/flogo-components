# MQTT Topic Subscriber
This trigger provides your flogo application the ability to start a flow via MQTT
It is different from the original one by Michael Register <mregiste@tibco.com> in the sense that it takes wildcards per endpoint and returns the actual topic that is used in a separate output.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/mqtt
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/mqtt
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
    {
      "name": "broker",
      "type": "string"
    },
    {
      "name": "id",
      "type": "string"
    },
    {
      "name": "user",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "actualtopic",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "topic",
        "type": "string",
      },
      {
        "name": "qos",
        "type": "integer",
        "allowed" : [0, 1, 2]
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| broker    | the MQTT Broker URI (tcp://[hostname]:[port])|
| id        | The MQTT Client ID |         
| user      | The UserID used when connecting to the MQTT broker |
| password  | The Password used when connecting to the MQTT broker |


## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message    | The message payload |
| actualtopic | The actual topic that was used to publish the message on) |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| topic    | The trigger will subscribe to this topic. May contain wildcards |
| qos      | MQTT Quality of Service. 0 = At most once, 1 = At least once, 2 = Exactly once. |

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the MQTT Trigger.

### Start a flow
Configure the Trigger to start "myflow". In this case the "endpoints" "settings" "topic" is "flogo/#" will start "testFlow" flow when a message arrives on a topic staring with "flogo" in this case. The actualtopic output will hold the actual topic used for further processing.

```json
{
  "name": "mqtt",
  "settings": {
    "broker": "tcp://127.0.0.1:1883",
    "id": "flogoEngine",
    "user": "",
    "password": ""
  },
  "handlers": [
    {
      "action":{
		"id":"dummy"
	  },
      "settings": {
        "topic": "flogo/#",
		"qos": "0"
      }
    }
  ]
}
```
