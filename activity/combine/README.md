# Combine
This activity provides your flogo application the ability to combine separate parts into a single string


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/combine
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/combine
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "delimiter",
      "type": "string"
    },
    {
      "name": "prefix",
      "type": "string"
    },
    {
      "name": "suffix",
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
## Inputs
| Input   | Description    |
|:----------|:---------------|
| delimiter    | The delimiter string, e.g."/" |
| prefix    | The fixed beginning of the result string  |
| suffix    | The fixed ending of the result string  |
| part1    | The value of the 1st part of the string |
| part2    | The value of the 2nd part of the string |
| part3    | The value of the 3rd part of the string |
| part4    | The value of the 4th part of the string |
| part5    | The value of the 5th part of the string |
| part6    | The value of the 6th part of the string |
| part7    | The value of the 7th part of the string |
| part8    | The value of the 8th part of the string |

## Outputs
| Output   | Description    |
|:----------|:---------------|
| result    | The resulting string output |

## Configuration Examples
### Simple
Configure a task in flow to combine 3 parts into one string 'flogo/device/DEV12345', using '/' as delimiter":

```json
{
  "name": "testcombine",
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
        "name": "Combine",
        "type": 1,
        "activityType": "combine",
        "attributes": [
          {
            "name": "delimiter",
            "value": "/",
            "type": "string"
          },
          {
            "name": "part1",
            "value": "flogo",
            "type": "string"
          },
          {
            "name": "part2",
            "value": "device",
            "type": "string"
          },
          {
            "name": "part3",
            "value": "DEV12345",
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
