![gofmt status](https://img.shields.io/badge/gofmt-compliant-green.svg?style=flat-square) ![golint status](https://img.shields.io/badge/golint-compliant-green.svg?style=flat-square) ![automated test coverage](https://img.shields.io/badge/test%20coverage-1%20testcase-orange.svg?style=flat-square)

# Send eFTL Message
This activity sends a message to TIBCO eFTL.
Code is based on the WSMessage activity created by Leon Stigter


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/eftl
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/eftl
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
      "name": "clientid",
      "type": "string",
      "value": ""
    },
    {
      "name": "channel",
      "type": "string",
      "value": ""
    },
    {
      "name": "destination",
      "type": "string",
      "value": ""
    },
    {
      "name": "subject",
      "type": "string",
      "value": ""
    },
    {
      "name": "message",
      "type": "string",
      "value": ""
    },
    {
      "name": "user",
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
  "output": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| server      | The WebSocket server to connect to (e.g. `localhost:9191`) |         
| clientid    | A client id used to identify the connection to eFTL |         
| channel     | The channel to send the message to (e.g. `/channel`)   |
| destination | The destination to send the message to (e.g. `default`) |
| subject     | The subject to pinpoint the message context (e.g. `sensor1`) |
| message     | The actual message to send |
| username    | The username to connect to the WebSocket server (e.g. `user`) |
| password    | The password to connect to the WebSocket server (e.g. `password`) |
| secure      | Determines to use secure web sockets (wss) |
| cert        | The eFTL server certificate data in base64 encoded PEM format (only used if secure is true) |

## Configuration Examples
The below configuration would connect to a WebSocket server based on TIBCO eFTL and send a message saying `Hello World`
```json
      {
        "id": 2,
        "name": "Send a message to eFTL",
        "type": 1,
        "activityType": "eftl",
        "attributes": [
          {
            "name": "server",
            "value": "localhost:9191",
            "type": "string"
          },
          {
            "name": "clientid",
            "value": "flogo_app",
            "type": "string"
          },
          {
            "name": "channel",
            "value": "/channel",
            "type": "string"
          },
          {
            "name": "destination",
            "value": "default",
            "type": "string"
          },
          {
            "name": "message",
            "value": "Hello World",
            "type": "string"
          },
          {
            "name": "user",
            "value": "user",
            "type": "string"
          },
          {
            "name": "password",
            "value": "password",
            "type": "string"
          }
        ]
      }
```

## Contributors
[Leon Stigter](https://github.com/retgits)
[Jan van der Lugt](https://github.com/jvanderl)
