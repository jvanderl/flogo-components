# Trim
This activity provides your flogo application the ability to trim data up to 8 strings

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/trim
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/trim
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "find",
      "type": "string"
    },
    {
      "name": "replace",
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
  "outputs": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "output1",
      "type": "string"
    },
    {
      "name": "output2",
      "type": "string"
    },
    {
      "name": "output3",
      "type": "string"
    },
    {
      "name": "output4",
      "type": "string"
    },
    {
      "name": "output5",
      "type": "string"
    },
    {
      "name": "output6",
      "type": "string"
    },
    {
      "name": "output7",
      "type": "string"
    },
    {
      "name": "output8",
      "type": "string"
    }
  ]
}
```
## Inputs
| Input   | Description    |
|:----------|:---------------|
| find     | What to replace from the strings stated below |
| replace  | What to replace the previous value by |
| input1   | The 1st string to trim |
| input2   | The 2nd string to trim |
| input3   | The 3rd string to trim |
| input4   | The 4th string to trim |
| input5   | The 5th string to trim |
| input6   | The 6th string to trim |
| input7   | The 7th string to trim |
| input8   | The 8th string to trim |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| result    | Result of the operation. Returns "OK" when trim succeeded |
| output1   | The trimmed version of input 1 |
| output2   | The trimmed version of input 2 |
| output3   | The trimmed version of input 3 |
| output4   | The trimmed version of input 4 |
| output5   | The trimmed version of input 5 |
| output6   | The trimmed version of input 6 |
| output7   | The trimmed version of input 7 |
| output8   | The trimmed version of input 8 |


## Configuration Examples
### Simple
Configure a task in flow to trim the colons from a mac address":

```json
{
  "name": "testtrim",
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
        "name": "Replace",
        "type": 1,
        "activityType": "replace",
        "attributes": [
          {
            "name": "find",
            "value": ":",
            "type": "string"
          },
          {
            "name": "replace",
            "value": "",
            "type": "string"
          },
          {
            "name": "input1",
            "value": "5C:CF:7F:94:2B:CB",
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
            "value": "{A2.outpu1}",
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
