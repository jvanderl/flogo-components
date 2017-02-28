# eFTL Subscriber
This trigger provides your flogo application the ability to start a flow via eFTL


## Installation

```bash
flogo add trigger github.com/jvanderl/flogo-components/trigger/eftl
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
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
      "name": "username",
      "type": "string",
      "value": ""
    },
    {
      "name": "password",
      "type": "string",
      "value": ""
    },
    {
      "name": "secure",
      "type": "boolean",
      "value": "false"
    },
    {
      "name": "certificate",
      "type": "string",
      "value": ""
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
        "name": "destination",
        "type": "string"
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| server    | the eFTL server [hostname]:[port]|
| channel     | The channel to send the message to (e.g. `/channel`)   |
| message     | The actual message to send |
| user        | The username to connect to the WebSocket server (e.g. `user`) |
| password    | The password to connect to the WebSocket server (e.g. `password`) |
| secure      | Determines to use secure web sockets (wss) |
| cert        | The eFTL server certificate data in PEM format (only used if secure is true) |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message    | The message payload |

## Endpoints
| Endpoint   | Description    |
|:----------|:---------------|
| destination | The destination to send the message to (e.g. `default`) |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the eFTL Trigger.

### Start a flow
Configure the Trigger to start "myflow". So in this case the "endpoints" "settings" "destination" is "flogo" will start "myflow" flow when a message arrives on a destination called "flogo" in this case.

```json
{
  "name": "eftl",
  "settings": {
    "server": "192.168.178.41:9191",
    "channel": "/channel",
    "user": "user",
    "password": "password"
  },
  "endpoints": [
    {
      "actionType": "flow",
      "actionURI": "local://testFlow",
      "settings": {
        "destination": "flogo"
      }
    }
  ]
}
```
