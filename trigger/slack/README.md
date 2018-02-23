# Slack Subscriber
This trigger provides your flogo application the ability to start a flow via a Slack message


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/slack
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/slack
```

## Schema
Settings, Outputs and Handler:

```json
{
  "settings":[
    {
      "name": "token",
      "type": "string",
      "required" : true
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "channel",
      "type": "string"
    },
    {
      "name": "username",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "channel",
        "type": "string",
        "required" : true
      },
      {
        "name": "matchtext",
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
| token      | Bot Token to connect to slack. Check 'Setup Slack' section here: https://rsmitty.github.io/Slack-Bot/  |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| message   | The message text |
| channel   | The channel the message was received on |
| username  | The user that sent the message |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| channel  | The channel to listen on. Use '*' for all channels |
| matchtext | Only trigger on messages containing this text. Use '*' for any text |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the TCM Trigger.

### Start a flow
Configure the Trigger to start "testFlow". So in this case any message received on channel "flogo" will start "testFlow".

```json
{
  "name": "slacklistener",
  "settings": {
    "token": "<your token here>"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "channel": "flogo",
        "matchtext": "*"
      }
    }
  ]
}
```
