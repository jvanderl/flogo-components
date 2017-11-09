# JSON to Data
This activity provides your flogo application the ability to convert a JSON string into a data object


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/jsontodata
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/jsontodata
```

## Schema
Inputs and Outputs:

```json
{
 "inputs":[
    {
      "name": "input",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "string"
    },    {
      "name": "data",
      "type": "object"
    }
  ]
}
```
## Inputs
| Input   | Description    |
|:----------|:---------------|
| input    | the input data offered as JSON string |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| result    | Result of the operation |
| data    | Data object representation of the input string |


## Configuration Examples
### Simple
Configure a task in flow to split a JSON string into separate output parameters":

```json
{
  "name": "testconverter",
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
        "name": "JSON to Data",
        "type": 1,
        "activityType": "jsontodata",
        "attributes": [
          {
            "name": "input",
            "value": "{\"distance\":1, \"deviceid\":\"DEV12345\"}",
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
            "value": "{A2.data}",
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
