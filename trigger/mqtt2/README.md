# MQTT Topic Subscriber
This trigger provides your flogo application the ability to start a flow via MQTT
It is different from the original one by Michael Register <mregiste@tibco.com> in the sense that it takes wildcards per endpoint and returns the actual topic that is used in a spearate output.


## Installation

```bash
flogo add trigger github.com/jvanderl/flogo-components/trigger/mqtt2
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
    },
    {
      "name": "store",
      "type": "string"
    },
    {
      "name": "topic",
      "type": "string"
    },
    {
      "name": "qos",
      "type": "number"
    },
    {
      "name": "cleansess",
      "type": "boolean"
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "actualtopic",
      "type": "string"
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "topic",
        "type": "string"
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
| store     | MQTT store for message persistence when QoS=1 or QoS=2
| topic     | Topic on which the message is published |
| qos       | MQTT Quality of Service |
| cleansess | Determines if the trigger should start with a clean session |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message    | The message payload |
| actualtopic | The actual topic that was used to publish the message on) |

## Endpoints
| Endpoint   | Description    |
|:----------|:---------------|
| topic    | The trigger will subscribe to this topic. May contain wildcards |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the MQTT Trigger.

### Start a flow
Configure the Trigger to start "myflow". "settings" "topic" is not used. So in this case the "endpoints" "settings" "topic" is "flogo/#" will start "myflow" flow when a message arrives on a topic staring with "flogo" in this case. The actualtopic output will hold the actucal topic used for further processing. 

```json
{
  "triggers": [
    {
      "name": "mqtt2",
      "settings": {
        "topic": "",
        "broker": "tcp://192.168.1.12:1883",
        "id": "flogo",
        "user": "",
        "password": "",
        "store": "",
        "qos": "0",
        "cleansess": "false"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "topic": "flogo/#"
          }
        }
      ]
    }
  ]
}
```
