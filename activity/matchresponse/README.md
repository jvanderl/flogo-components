# Match Response
This activity provides your flogo application the ability to find a substring and send a response based on provided search data. One use case is for a bot.

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/matchresponse
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/matchresponse
```

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
      "name": "searchdata",
      "type": "array"
    }
  ],
  "outputs": [
    {
      "name": "match",
      "type": "string"
    },
    {
      "name": "response",
      "type": "string"
    }
  ]
}
```
## Inputs
| Input   | Description    |
|:----------|:---------------|
| input     | The string to find a substring in |
| searchdata  | array of search data. Format is [{"find":"lookuptxt1","resp":"response1"},{"find":"lookuptxt2","resp":"response2"}] |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| match    | The match found in "find" element, "not found" if no match |
| respsonse   | The response that comes with the matched lookup. "not found" if no match |

## Configuration Examples
### Simple
Configure a task in flow to send response 'Don't mention it' when input contains 'thank you'":

```json
{
  "name": "testmatchresponse",
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
        "name": "MatchResponse",
        "type": 1,
        "activityType": "matchresponse",
        "attributes": [
          {
            "name": "input",
            "value": "thank you for the music!",
            "type": "string"
          },
          {
            "name": "searchdata",
            "value": "[{\"find\": \"thank you\", \"resp\": \"Don't mention it\"},{\"find\": \"bye\", \"resp\": \"See you later!\"}
		]`",
            "type": "array"
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
            "value": "{A2.reponse}",
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
