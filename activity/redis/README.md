# Interact with Redis
This activity provides your flogo application the ability to interact with a Redis keyspace.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/redis
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/redis
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
   {
      "name": "server",
      "type": "string",
      "value": "localhost:6379"
    },
    {
      "name": "password",
      "type": "string"
    },
    {
      "name": "database",
      "type": "integer"
    },
    {
      "name": "operation",
      "type": "string",
      "required": true,
      "allowed" : ["GET", "PUT", "DEL", "PING"]
    },
    {
      "name": "key",
      "type": "string"
    },
    {
      "name": "value",
      "type": "string"
    }
  ],
  "outputs": [
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
| server    | the address of the redis server ([hostname]:[port])|
| password  | The Password used when connecting to the redis server |
| database  | Redis database, default 0 |
| operation | can be "GET", "SET", "DEL", "PING" |
| key       | Key used to interact with redis |
| value     | The value to write (only applies when operation is "SET") |


## Configuration Examples
### Simple
Configure a task in flow to write "hello world" to a key on redis called "flogo":

```json
{
  "id": 2,
  "name": "Interact with Redis",
  "type": 1,
  "activityType": "redis",
  "attributes": [
    {
      "name": "server",
      "value": "localhost:6379",
      "type": "string"
    },
    {
      "name": "password",
      "value": "",
      "type": "string"
    },
    {
      "name": "database",
      "value": "0",
      "type": "integer"
    },
    {
      "name": "operation",
      "value": "GET",
      "type": "string"
    },
    {
      "name": "key",
      "value": "flogo",
      "type": "string"
    },
    {
      "name": "value",
      "value": "Hello World",
      "type": "string"
    }
  ]
}
```
