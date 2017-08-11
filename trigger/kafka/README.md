# Kafka Topic Subscriber
This trigger provides your flogo application the ability to start a flow via Kafka


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/kafka
```
Link for flogo web: https://github.com/jvanderl/flogo-components/trigger/kafka

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
    {
      "name": "server",
      "type": "string"
    },
    {
      "name": "configid",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "topic",
        "type": "string"
      },
      {
        "name": "partition",
        "type": "int"
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| server    | the Kafka server(s) [hostname]:[port]|
| configid       | The Kafka client configuration ID |         

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message    | The message payload |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| topic    | The trigger will subscribe to this topic. |
| partition  | Partition to use for the subscription |

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Kafka Trigger.

### Start a flow
Configure the Trigger to start "myflow". So in this case the "handlers" "settings" "topic" is "flogo" will start "testFlow" flow when a message arrives on a topic staring with "flogo" in this case.

```json
{
  "name": "kafka",
  "settings": {
    "server": "127.0.0.1:9092",
    "configid": "test-flogo-trigger"
  },
  "handlers": [
		{
      "actionId": "local://testFlow",
      "settings": {
        "topic": "flogo",
				"partition": "0"
      }
    }
  ]
}

```
