# udp
This activity provides your flogo application the ability to send a message over UDP.


## Installation

```bash
flogo add activity github.com/jvanderl/flogo-components/incubator/activity/udp
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "port",
      "type": "integer",
      "required": true
    },
    {
      "name": "multicast_group",
      "type": "string",
      "required": true
    },
    {
      "name": "message",
      "type": "string",
      "required": true
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "any"
    }
  ]
}

```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| port      | The UDP portnumber |         
| multicast_group         | The IP address of the multicast group   |
| message  | The message to send |

## Outputs
| Output     | Description    |
|:------------|:---------------|
| result      | 'OK' when message is sent correctly, otherwise error result |         

## Configuration Examples
### Simple
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io):

```json
{
  "id": 3,
  "type": 1,
  "activityType": "udp",
  "name": "Send Hello World",
  "attributes": [
    { "name": "port", "value": "10001" },
    { "name": "multicast_group", "value": "127.0.0.1" },
    { "name": "message", "value": "Hello World" }    
  ]
}
```
