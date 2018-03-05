# State Change
This activity provides your flogo application the ability to detect state change for up to eight inputs.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/statechange
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/statechange
```

## Schema
Inputs and Outputs:

```json
{
 "input":[
    {
      "name": "datasource",
      "type": "string"
    },
    {
      "name": "input1",
      "type": "string"
    },
    {
      "name": "input2",
      "type": "string"
    },
    {
      "name": "input3",
      "type": "string"
    },
    {
      "name": "input4",
      "type": "string"
    },
    {
      "name": "input5",
      "type": "string"
    },
    {
      "name": "input6",
      "type": "string"
    },
    {
      "name": "input7",
      "type": "string"
    },
    {
      "name": "input8",
      "type": "string"
    }   
  ],
  "output": [
    {
      "name": "changed",
      "type": "boolean"
    },
    {
      "name": "flags",
      "type": "integer"
    },
    {
      "name": "result",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| datasource | For identifying unique values |
| input 1-8  | input values for detecting change. Make sure you use the same inputs for each type of value across your flow|

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| changed | Determines if any of the input values have changed (true, false) |
| flags  | Bitwise representation wat input value has changed. flags & 1 = input 1, flags & 2 = intput 2, flags & 4 = input 3, flags & 8 = input 4, flags & 16 = input 5 etc. |
| result  | Result of state change in text, return "NO_CHANGE" if no input changed state |


## Configuration Examples
### Simple
Configure a task in flow to only pass data only once per 10 seconds:

```json
{
  "name": "teststatechange",
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
        "name": "State Change",
        "type": 1,
        "activityType": "statechange",
        "attributes": [
          {
            "name": "datasource",
            "value": "flogo/device/DEV12345/distance",
            "type": "string"
          },
          {
            "name": "input1",
            "value": "10",
            "type": "integer"
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
