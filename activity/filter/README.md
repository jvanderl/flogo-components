# Publish MQTT Message
This activity provides your flogo application the ability to filter out unwanted data.


## Installation

```bash
flogo add activity github.com/jvanderl/flogo-components/activity/filter
```

## Schema
Inputs and Outputs:

```json
{
 "inputs":[
    {
      "name": "input",
      "type": "string",
      "required": true
    },
    {
      "name": "datasource",
      "type": "string"
    },
    {
      "name": "datatype",
      "type": "string",
      "required": true,
      "allowed" : ["int", "uint", "float32"]
    },
    {
      "name": "minvalue",
      "type": "string"
    },
    {
      "name": "maxvalue",
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
      "name": "result",
      "type": "string"
    },
    {
      "name": "usevalue",
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
| input    | the input data |
| datasource | For identifying unique values (only used i.c.w. interval) |
| datatype  | The type of data offert (int, uint or float32) |
| minvalue  | The minimum value that gets passed through |
| maxvalue  | The maximum value that gets passed through |
| interval  | The time interval for passing unique data|
| intervaltype | Unit of time for interval (hours, minutes, seconds or milliseconds) |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| result    | The output data when not filtered out |
| usevalue | Determines if the result should be used or not (true, false) |
| reason  | When the data is not to be used, reason explains why |
| lasttimepassed  | When the filtering by interval, this holds the time when data was last used |


## Configuration Examples
### Simple
Configure a task in flow to only forward data between 100 and 200 with a 10 second interval:

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
        "name": "Filter",
        "type": 1,
        "activityType": "filter",
        "attributes": [
          {
            "name": "input",
            "value": "125",
            "type": "string"
          },
          {
            "name": "datasource",
            "value": "flogo/device/DEV12345/distance",
            "type": "string"
          },
          {
            "name": "datatype",
            "value": "int",
            "type": "string"
          },
          {
            "name": "minvalue",
            "value": "100",
            "type": "string"
          },
          {
            "name": "maxvalue",
            "value": "200",
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
            "value": "{A2.result}",
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
