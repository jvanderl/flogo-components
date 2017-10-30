# [title]

This activity provides your flogo application the ability to [functionality]

## Installation

```bash
flogo install github.com/[git_user]/[git_repo]/activity/[package]
```
Link for flogo web:
```
https://github.com/[git_user]/[git_repo]/activity/[package]
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "[input1_name]",
      "type": "[input1_type]"
    },
    {
      "name": "[input2_name]",
      "type": "[input2_type]"
    }
  ],
  "outputs": [
    {
      "name": "[output1_name]",
      "type": "[output1_type]"
    },
    {
      "name": "[output2_name]",
      "type": "[output2_type]"
    }
  ]
}
```
## Inputs
| Input   | Description    |
|:----------|:---------------|
| [input1_name]     | [input1_desc] |
| [input2_name]     | [input2_desc] |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| [output1_name]     | [output1_desc] |
| [output2_name]     | [output2_desc] |


## Configuration Examples
### Simple
[example_desc]:

```json
{
  "name": "test[package]",
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
        "name": "[title]",
        "type": 1,
        "activityType": "[package]",
        "attributes": [
          {
            "name": "[input1_name]",
            "value": "[input1_value]",
            "type": "[input1_type]"
          },
          {
            "name": "[input2_name]",
            "value": "[input2_value]",
            "type": "[input2_type]"
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
            "value": "{A2.[output1_name]}",
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
