
# Send Slack Message
This activity sends a message to Slack.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/slack
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/slack
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "token",
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
      "name": "result",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting          | Description    |
|:-----------------|:---------------|
| token              | Bot Token to connect to slack. Check 'Setup Slack' section here: https://rsmitty.github.io/Slack-Bot/ |         
| channel         | The Slack channel to publish on |
| message         | The message to send |         

## Configuration Examples
The below configuration would connect to Slack and send a message saying `Hello World`
```json
      {
        "id": 2,
        "name": "Send a message to Slack",
        "type": 1,
        "activityType": "slackpub",
        "attributes": [
          {
            "name": "tokem",
            "value": "<your token here>",
            "type": "string"
          },
          {
            "name": "channel",
            "value": "flogo",
            "type": "string"
          },
          {
            "name": "message",
            "value": "Hello World",
            "type": "string"
          }
        ]
      }
```
