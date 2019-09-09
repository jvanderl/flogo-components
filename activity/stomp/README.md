
# Send Message on Stomp
This activity sends a message to Slack.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/stomp
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/stomp
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "address",
      "type": "string",
      "value": ""
    },
    {
      "name": "destination",
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
| Setting      | Description    |
|:-------------|:---------------|
| address      | Adress for stomp server. Example: localhost:61613 |         
| destination  | The destination to send message on |
| message      | The message to send |         

## Configuration Examples
The below configuration would connect to Stomp server on `locahost:61613` and send a message on destination `flogo` saying `Hello World`
```json
      {
        "id": 2,
        "name": "Send a message on Stomp",
        "type": 1,
        "activityType": "stomp",
        "attributes": [
          {
            "name": "address",
            "value": "localhost:61613",
            "type": "string"
          },
          {
            "name": "detination",
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
