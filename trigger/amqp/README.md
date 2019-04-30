# AMQP Subscriber
This trigger provides your flogo application the ability to start a flow via AMQP


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/amqp
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/amqp
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
    {
      "name": "server",
      "type": "string",
      "required": true,
      "value" : "localhost"
    },
    {
      "name": "port",
      "type": "string",
      "required": true,
      "value" : "5672"
    },
    {
      "name": "userID",
      "type": "string",
      "required": true,
      "value" : "guest"
    },
    {
      "name": "password",
      "type": "string",
      "required": true,
      "value" : "guest"
    }
  ],
  "output": [
    {
      "name": "message",
      "type": "any"
    },
    {
      "name": "contentType",
      "type": "string"
    },
    {
      "name": "routingKey",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "exchange",
        "type": "string"
      },
      {
        "name": "exchangeType",
        "type": "string",
        "allowed" : ["direct", "fanout", "topic", "x-custom"],
        "value" : "direct"
      },
      {
        "name": "queueName",
        "type": "string"
      },
      {
        "name": "bindingKey",
        "type": "string"
      },
      {
        "name": "consumerTag",
        "type": "string"
      },
      {
        "name": "durable",
        "type": "boolean",
        "value": "true"
      },
      {
        "name": "autoDelete",
        "type": "boolean",
        "value": "false"
      },
      {
        "name": "exclusive",
        "type": "boolean",
        "value": "false"
      },
      {
        "name": "noWait",
        "type": "boolean",
        "value": "false"
      },
      {
        "name": "noAck",
        "type": "boolean",
        "value": "false"
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| server      | The AMQP Host Name (default "localhost") |
| port        | The AMQP Port Number (default "5672") |         
| userID      | The UserID used when connecting to the AMQP server |
| password    | The Password used when connecting to the AMQP server |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message   | The message payload |
| contentType | The content type of the message |
| routingKey   | The original routing key used by the AMQP publisher (handy when listening to wildcards) |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| exchange | The exchange name used for the subscriber |
| exchangeType | The type of exchange. can be "direct", "fanout", "topic" or "x-custom" |
| queueName | The queue name to listen to |
| bindingKey | The binding key expresses interest. Use # for all, refer to AMQP documentation online |
| consumerTag | Free to assign tag for the message consumer |
| durable     | Creates durable topic / queue when set to true (default "true") |
| autoDelete  | Deletes topic / queue when no longer used (default "false") |
| exclusive   | Sets topic / queue up for internal use when set to true (default "false") |
| noWait      | When set to true the server will not respond to the message (default "false") |
| noAck      | When set to true the server will not acknowledge the received message (default "false") |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the AMQP Trigger.

### Start a flow
Configure the Trigger to start "testFlow". So in this case the "endpoints" "settings" "destination" is "flogo" will start "testFlow" flow when a message arrives on a queue called "flogo" in this case.

```json
{
  "name": "eftl",
  "settings": {
        "server": "localhost",
        "port": "5672",
        "userID": "guest",
        "password": "guest"
    },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
            "exchange": "test",
            "exchangeType": "direct",
            "queueName": "flogo",
            "bindingKey": "test",
            "consumerTag": "test",
            "durable": "true",
            "autoDelete": "false",
            "exclusive": "false",
            "noWait": "false",
            "noAck": "false"
      }
    }
  ]
}
```
