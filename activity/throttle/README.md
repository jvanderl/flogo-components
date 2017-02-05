# Throttle by Interval
This activity provides your flogo application the ability to throttle data by interval.


## Installation

```bash
flogo add activity github.com/jvanderl/flogo-components/activity/throttle
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "datasource",
      "type": "string"
    },
    {
      "name": "interval",
      "type": "integer"
    },
    {
      "name": "intervaltype",
      "type": "string",
      "allowed" : ["hours", "minutes", "seconds", "milliseconds"]
    }
  ],
  "outputs": [
    {
      "name": "pass",
      "type": "boolean"
    },
    {
      "name": "reason",
      "type": "string"
    },
    {
      "name": "lasttimepassed",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| datasource | For identifying unique values |
| interval  | The time interval for passing unique data |
| intervaltype | Unit of time for interval (hours, minutes, seconds or milliseconds) |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| pass | Determines if the data should be passed (true, false) |
| reason  | When the data is not to be passed, reason explains why |
| lasttimepassed  | This holds the time when data passed the throttle |


## Configuration Examples
### Simple
Configure a task in flow to only pass data only once per 10 seconds:

```json
{
  "name": "testfilter",
  "model": "tibco-simple",
  "type": 1,
  "attributes": [],
  "rootTask": {
    "id": 1,
    "type": 1,
    "activityType": "",
    "name": "root",
    "tasks": [
      {
        "id": 2,
        "name": "Throttle",
        "type": 1,
        "activityType": "throttle",
        "attributes": [
          {
            "name": "datasource",
            "value": "flogo/device/DEV12345/distance",
            "type": "string"
          },
          {
            "name": "interval",
            "value": "10",
            "type": "integer"
          },
          {
            "name": "intervaltype",
            "value": "seconds",
            "type": "string"
          }
        ]
      },
      {
        "id": 3,
        "name": "Log Message",
        "type": 1,
        "activityType": "tibco-log",
        "attributes": [
          {
            "name": "message",
            "value": "",
            "type": "string"
          },
          {
            "name": "flowInfo",
            "value": "true",
            "type": "boolean"
          },
          {
            "name": "addToFlow",
            "value": "true",
            "type": "boolean"
          }
        ],
        "inputMappings": [
          {
            "type": 1,
            "value": "{A2.reason}",
            "mapTo": "message"
          }
        ]
      }
    ],
    "links": [
      {
        "id": 1,
        "from": 2,
        "to": 3,
        "type": 0
      }
    ]
  }
}
```
