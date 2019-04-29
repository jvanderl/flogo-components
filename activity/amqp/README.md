# Publish AMQP Message
This activity provides your flogo application the ability to publish a message on an AMQP queue or topic.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/amqp
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/amqp
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
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
    },
    {
      "name": "message",
      "type": "any",
      "required": true
    },
    {
      "name": "exchange",
      "type": "string",
      "required": false,
      "value" : ""
    },
    {
      "name": "routingKey",
      "type": "string",
      "required": true
    },
    {
      "name": "routingType",
      "type": "string",
      "required": true,
      "allowed" : ["queue", "topic"]
    },
    {
      "name": "durable",
      "type": "boolean",
      "required": false,
      "value" : "false"
    },
    {
      "name": "autoDelete",
      "type": "boolean",
      "required": false,
      "value" : "false"
    },
    {
      "name": "exclusive",
      "type": "boolean",
      "required": false,
      "value" : "false"
    },
    {
      "name": "noWait",
      "type": "boolean",
      "required": true,
      "value" : "false"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "any"
    }
  ]
}
```
## Inputs
| Input       | Description    |
|:------------|:---------------|
| server      | The AMQP Host Name (default "localhost") |
| port        | The AMQP Port Number (default "5672") |         
| userID      | The UserID used when connecting to the AMQP server |
| password    | The Password used when connecting to the AMQP server |
| message     | The message payload |
| exchange    | AMPQ Exhcange name. (default "" when using queues) |
| routingKey  | AMPQ routing key when using topics. When using queues, put queue name here |
| routingType | Type of routing. Choose "queue" or "topic" (default "queue") |
| durable     | Creates durable topic / queue when set to true (default "false") |
| autoDelete  | Deletes topic / queue when no longer used (default "false") |
| exclusive   | Sets topic / queue up for internal use when set to true (default "false") |
| noWait      | When set to true the server will not respond to the message (default "false") |


## Configuration Examples
### Simple
Configure a task in flow to publish a "hello world" message on AMQP queue called "hello":

```json
{
  "id": 2,
  "name": "Publish AMQP Message",
  "type": 1,
  "activityType": "amqp",
  "attributes": [
    {
      "name": "server",
      "value": "localhost",
      "type": "string"
    },
    {
      "name": "port",
      "value": "5672",
      "type": "number"
    },
    {
      "name": "userID",
      "value": "guest",
      "type": "string"
    },
    {
      "name": "password",
      "value": "guest",
      "type": "string"
    },
    {
      "name": "message",
      "value": "Hello World, from Flogo!",
      "type": "string"
    },
    {
      "name": "routingKey",
      "value": "hello",
      "type": "string"
    }
  ]
}
```
