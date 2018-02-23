
# Send TCM Message
This activity sends a message to TIBCO Cloud Messaging.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/tcmpub
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/tcmpub
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "url",
      "type": "string",
      "value": ""
    },
    {
      "name": "authkey",
      "type": "string",
      "value": ""
    },
    {
      "name": "clientid",
      "type": "string",
      "value": ""
    },
    {
      "name": "destinationname",
      "type": "string",
      "value": ""
    },
    {
      "name": "destinationvalue",
      "type": "string",
      "value": ""
    },
    {
      "name": "messagename",
      "type": "string",
      "value": ""
    },
    {
      "name": "messagevalue",
      "type": "string",
      "value": ""
    },
    {
      "name": "certificate",
      "type": "string",
      "value": ""
    }
  ],
  "outputs": [
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
| url              | The TIBCO Cloud Messaging URL (wss://nn.messaging.cloud.tibco.com/tcm/xxxxx/channel ) |         
| authkey         | The TIBCO Authorization Key |
| clientid         | A unique client id used to identify the connection to TCM |         
| destinationname  | The identifier of the destination field (can be left empty) |
| destinationvalue | The destination to send the message on (can be left empty) |
| messagename      | The identifier of the message field |
| messagevalue     | The actual message to send |
| cert             | The TCM server certificate data in base64 encoded PEM format (accepts all when empty) |

## Configuration Examples
The below configuration would connect to TIBCO Cloud Messaging and send a message saying `Hello World`
```json
      {
        "id": 2,
        "name": "Send a message to TCM",
        "type": 1,
        "activityType": "tcmpub",
        "attributes": [
          {
            "name": "url",
            "value": "wss://nn.messaging.cloud.tibco.com/tcm/xxxxx/channel",
            "type": "string"
          },
          {
            "name": "authkey",
            "value": "XYZXYZXYZXYZXYZXYZ",
            "type": "string"
          },
          {
            "name": "clientid",
            "value": "flogo_app",
            "type": "string"
          },
          {
            "name": "destinationname",
            "value": "",
            "type": "string"
          },
          {
            "name": "destinationvalue",
            "value": "",
            "type": "string"
          },
          {
            "name": "messagename",
            "value": "demo_tcm",
            "type": "string"
          },
          {
            "name": "messagevalue",
            "value": "Hello World",
            "type": "string"
          }
        ]
      }
```

## Contributors
[Leon Stigter](https://github.com/retgits)

[Jan van der Lugt](https://github.com/jvanderl)
