# WsServer
This trigger provides your flogo application the ability to start a flow when you receive a Websocket Message.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/incubator/trigger/wsserver
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/incubator/trigger/wsserver
```

## Schema
Outputs and Endpoint:

```json
{
  "settings": [
  ],
  "output": [
        {
      "name": "params",
      "type": "params"
    },
    {
      "name": "triggerTime",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
    ]
  }
}
```
## Settings
- None -

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| triggerTime |  The date and time the trigger fired |

## Handlers
- None -

## Example Configuration

Trigger to run a flow at startup

```json
{
  "name": "onstartup",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
      }
    }
  ]
}
```