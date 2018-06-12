# Send WebSocket message
This activity allows you to send a message to a WebSocket server.

## Installation
```bash
flogo add activity github.com/jvanderl/flogo-components/activity/wssend
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/wssend
```

## Schema
Inputs and Outputs:

```json
{
"input":[
    {
      "name": "server",
      "type": "string",
      "value": ""
    },
    {
      "name": "channel",
      "type": "string",
      "value": ""
    },
    {
      "name": "message",
      "type": "string",
      "value": ""
    }
  ],
  "output": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| server      | False    | The WebSocket server to connect to (e.g. `localhost:9191`) |         
| channel     | False    | The channel to send the message to (e.g. `/channel`)   |
| message     | False    | The message to send |
## Outputs
| Output       | Description    |
|:------------|:---------------|
| output      | A string with the result of the action |

## Configuration Examples
The below example sends a message `Hello World`
```json
{
  "id": "wssend",
  "name": "Send WebSocket Message",
  "description": "This activity sends a message to a WebSocket server",
  "activity": {
    "ref": "github.com/jvanderl/flogo-components/activity/wssend",
    "input": {
      "Server": "localhost:9191",
      "Channel": "/channel",
      "Message": "Hello World"
    }
  }
}
```