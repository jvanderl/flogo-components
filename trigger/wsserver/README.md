# WsServer
This trigger provides your flogo application the ability to start a flow when you receive a Websocket Message.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/wsserver
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/wsserver
```

## Schema
Outputs and Endpoint:

```json
{
  "settings": [
    {
      "name": "port",
      "type": "integer"
    }
  ],
  "output": [
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "channel",
      "type": "string"
    }
  ],
  "reply": [
    {
      "name": "response",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
        {
          "name": "channel",
          "type": "string"
        }
      ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| port |  The port the built-in WebSocket Server is listening to |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message |  The message received by the trigger |
| channel |  The channel the message is received on |

## Reply
| Name   | Description    |
|:----------|:---------------|
| response |  The response message returned to the WebSocket sender |

## Handler Settings
| Setting   | Description    |
|:----------|:---------------|
| channel |  The channel messages can be received |

## Example Configuration

Trigger to run a flow when a message is receved on WS channel 'test'

```json
{
  "name": "wsserver",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "channel": "/test",
      }
    }
  ]
}
```