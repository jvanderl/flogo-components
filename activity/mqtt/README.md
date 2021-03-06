# Publish MQTT Message
This activity provides your flogo application the ability to publish a message on an MQTT topic.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/mqtt
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/mqtt
```

## Schema
Inputs and Outputs:

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
  "input":[

    {
      "name": "topic",
      "type": "string"
    },
    {
      "name": "qos",
      "type": "integer",
      "allowed" : [0, 1, 2]
    },
    {
      "name": "message",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| broker    | The MQTT Broker URI (tcp://[hostname]:[port])|
| id        | The MQTT Connection Client ID |         
| user      | The UserID used when connecting to the MQTT broker |
| password  | The Password used when connecting to the MQTT broker |

## Inputs
| Input   | Description    |
| topic     | Topic on which the message is published |
| qos       | MQTT Quality of Service. 0 = At most once, 1 = At least once, 2 = Exactly once. |
| message   | The message payload |


## Configuration Examples
### Simple
Configure a task in flow to publish a "hello world" message on MQTT topic called "flogo":

```json
{
  "id": "mqtt_3",
  "name": "Send MQTT Message",
  "description": "Pubishes message on MQTT topic",
  "activity": {
    "ref": "#mqtt",
    "input": {
      "message": "hello, world",
      "qos": 0,
      "topic": "flogo"
    },
    "settings": {
      "broker": "tcp://127.0.0.1:1883",
      "id": "flogo-test"
    }
  }
}
```
