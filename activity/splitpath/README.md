# Split Path
This activity provides your flogo application the ability to split a path into separate parts


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/splitpath
```
Link for flogo web: https://github.com/jvanderl/flogo-components/activity/splitpath

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "input",
      "type": "string"
    },
    {
      "name": "delimiter",
      "type": "string"
    },
    {
      "name": "fixedpath",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "fixedpath",
      "type": "string"
    },
    {
      "name": "part1",
      "type": "string"
    },
    {
      "name": "part2",
      "type": "string"
    },
    {
      "name": "part3",
      "type": "string"
    },
    {
      "name": "part4",
      "type": "string"
    },
    {
      "name": "part5",
      "type": "string"
    },
    {
      "name": "part6",
      "type": "string"
    },
    {
      "name": "part7",
      "type": "string"
    },
    {
      "name": "part8",
      "type": "string"
    },
    {
      "name": "remainder",
      "type": "string"
    }
  ]
}
```
## Inputs
| Input   | Description    |
|:----------|:---------------|
| input    | The input data offered as string some kind of delimiter |
| delimiter    | The delimiter string, e.g."/" |
| fixedpath    | The beginning of the string that is not to be split |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| result    | The result of the activity, "OK" if no errors |
| fixedpath  | The fixed part of the path copied over from the input |
| part1    | The value of the 1st part of the path, starting from the right |
| part2    | The value of the 2nd part of the path, starting from the right |
| part3    | The value of the 3rd part of the path, starting from the right |
| part4    | The value of the 4th part of the path, starting from the right |
| part5    | The value of the 5th part of the path, starting from the right |
| part6    | The value of the 6th part of the path, starting from the right |
| part7    | The value of the 7th part of the path, starting from the right |
| part8    | The value of the 8th part of the path, starting from the right |
| remainder | When ther are more than 8 parts, the rest of the path here |

## Configuration Examples
### Simple
Configure a task in flow to split a path into separate output parts":

```json
{
  "name": "testpath",
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
        "name": "Split Path",
        "type": 1,
        "activityType": "splitpath",
        "attributes": [
          {
            "name": "input",
            "value": "prefix/0/1/2/3/4/5/6/7/8/9",
            "type": "string"
          },
          {
            "name": "delimiter",
            "value": "/",
            "type": "string"
          },
          {
            "name": "fixedpath",
            "value": "prefix",
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
