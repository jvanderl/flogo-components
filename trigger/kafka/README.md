# Kafka Topic Subscriber
This trigger provides your flogo application the ability to start a flow via Kafka


## Installation

```bash
flogo add trigger github.com/jvanderl/flogo-components/trigger/kafka
```

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
    },
    {
      "name": "topic",
      "type": "string"
    },
    {
      "name": "partition",
      "type": "int"
    }
  ],
  "outputs": [
    {
      "name": "message",
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
| server    | the Kafka server(s) [hostname]:[port]|
| configid       | The Kafka client configuration ID |         
| topic     | Topic on which the message is published |
| partition  | Partition to use |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message    | The message payload |

## Endpoints
| Endpoint   | Description    |
|:----------|:---------------|
| topic    | The trigger will subscribe to this topic. |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Kafka Trigger.

### Start a flow
Configure the Trigger to start "myflow". "settings" "topic" is not used. So in this case the "endpoints" "settings" "topic" is "flogo/#" will start "myflow" flow when a message arrives on a topic staring with "flogo" in this case. The actualtopic output will hold the actucal topic used for further processing. 

```json
{
  "triggers": [
    {
      "name": "kafka",
      "settings": {
        "server": "tcp://192.168.1.12:1883",
        "configid": "flogo-test",
        "topic": "test",
        "partition": "0"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "topic": "test"
          }
        }
      ]
    }
  ]
}
```
