# eFTL Subscriber
This trigger provides your flogo application the ability to start a flow via eFTL


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/eftl
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/eftl
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
    {
      "name": "server",
      "type": "string",
      "required" : true
    },
    {
      "name": "clientid",
      "type": "string",
      "required" : true
    },
    {
      "name": "channel",
      "type": "string",
      "required" : true
    },
    {
      "name": "user",
      "type": "string",
      "required" : false
    },
    {
      "name": "password",
      "type": "string",
      "required" : false
    },
    {
      "name": "secure",
      "type": "boolean",
      "required" : true
    },
    {
      "name": "certificate",
      "type": "string",
      "required" : false
    }
  ],
  "output": [
    {
      "name": "message",
      "type": "object"
    }
  ],
  "reply": [
    {
      "name": "data",
      "type": "object"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "matcher",
        "type": "string",
        "required" : true
      },
      {
        "name": "durable",
        "type": "boolean",
        "required" : true
      },
      {
        "name": "durablename",
        "type": "string",
        "required" : false
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| server    | the eFTL server [hostname]:[port]|
| clientid    | the client id to identify the eFTL connection |
| channel     | The channel to send the message to (e.g. `/channel`)   |
| message     | The actual message to send |
| user        | The username to connect to the WebSocket server (e.g. `user`) |
| password    | The password to connect to the WebSocket server (e.g. `password`) |
| secure      | Determines to use secure web sockets (wss) |
| cert        | The eFTL server certificate data in base64 encoded PEM format (only used if secure is true) |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message   | The message object |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| destination | The destination to subscribe to (e.g. 'Default'), can also do '{"_dest":"subject"}' |
| usesubject | Use subject (set to false if you're only subscribing to a destination) |
| subject | Use this to match a specific subject of interest (when provided by sender) |
| durable | Use durable subscription |
| durablename | Durable subscription name (only needed when durable is set to 'true') |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the eFTL Trigger.

### Start a flow
Configure the Trigger to start "testFlow". So in this case the "endpoints" "settings" "destination" is "flogo" will start "testFlow" flow when a message arrives on a destination called "flogo" in this case.

```json
{
  "name": "eftl",
  "settings": {
    "server": "localhost:9191",
    "cclientid": "flogo-subscriber",
    "channel": "/channel",
    "user": "",
    "password": "",
    "secure": "false",
    "certificate": "*** Base64 encoded PEM cert data here ***"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "destination": "flogo",
        "subject": "sensor1",
        "durable": "false"
      }
    }
  ]
}
```
