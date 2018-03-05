# Split JSON
This activity provides your flogo application the ability to split a simple JSON structure into separate name-value pairs


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/splitjson
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/splitjson
```

## Schema
Inputs and Outputs:

```json
{
 "input":[
    {
      "name": "input",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    },    {
      "name": "name1",
      "type": "string"
    },
    {
      "name": "value1",
      "type": "string"
    },
    {
      "name": "name2",
      "type": "string"
    },
    {
      "name": "value2",
      "type": "string"
    },
    {
      "name": "name3",
      "type": "string"
    },
    {
      "name": "value3",
      "type": "string"
    },
    {
      "name": "name4",
      "type": "string"
    },
    {
      "name": "value4",
      "type": "string"
    },
    {
      "name": "name5",
      "type": "string"
    },
    {
      "name": "value5",
      "type": "string"
    },
    {
      "name": "name6",
      "type": "string"
    },
    {
      "name": "value6",
      "type": "string"
    },
    {
      "name": "name7",
      "type": "string"
    },
    {
      "name": "value7",
      "type": "string"
    },
    {
      "name": "name8",
      "type": "string"
    },
    {
      "name": "value8",
      "type": "string"
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
| name1    | The name of the 1st element in the JSON string |
| value1    | The value of the 1st element in the JSON string |
| name2    | The name of the 2nd element in the JSON string |
| value2    | The value of the 2nd element in the JSON string |
| name3    | The name of the 3rd element in the JSON string |
| value3    | The value of the 3rd element in the JSON string |
| name4    | The name of the 4th element in the JSON string |
| value4    | The value of the 4th element in the JSON string |
| name5    | The name of the 5th element in the JSON string |
| value5    | The value of the 5th element in the JSON string |
| name6    | The name of the 6th element in the JSON string |
| value6    | The value of the 6th element in the JSON string |
| name7    | The name of the 7th element in the JSON string |
| value7    | The value of the 7th element in the JSON string |
| name8    | The name of the 8th element in the JSON string |
| value8    | The value of the 8th element in the JSON string |


## Configuration Examples
### Simple
Configure a task in flow to split a JSON string into separate output parameters":

```json
{
  "name": "testsplitter",
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
        "name": "Split JSON",
        "type": 1,
        "activityType": "splitjson",
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
            "value": "{A2.name1}",
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
