# Stomp
This trigger provides your flogo application the ability to subscribe to messages on Stomp.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/stomp
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/stomp
```

## Schema
Outputs and Endpoint:

```json
{
  "settings":[
    {
      "name": "address",
      "type": "string",
      "description": "The address of the Stomp server to connect to"
    },
    {
      "name": "username",
      "type": "string",
      "description": "The username used to login to the Stomp server"
    },
    {
      "name": "password",
      "type": "string",
      "description": "The password used to login to the Stomp server"
    }
  ],
  "output": [
    {
      "name": "message",
      "type": "any"
    },
    {
      "name": "originalSource",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "source",
        "type": "string"
      }
    ]
  }
}

```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| address |  The address to connect to Stomp server. Example: localhost:61613 |
| username |  The username used to login to the Stomp server. Example admin |
| password |  The password used to login to the Stomp server. Example admin |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message |  The message that was received |
| originalSource |  The actual topic or queue the message was received from |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| source    | The topic or queue that the handler subscribes to |

## Example Configuration

Subscribe to messages on source `flogo` on stomp server `localhost:61613`

```json
{
  "name": "stomp",
  "settings": {
    "address": "localhost:61613",
    "username": "admin",
    "password": "admin",
  },
  "handlers": [
    {
      "source": "flogo",
    }
  ]
}
```