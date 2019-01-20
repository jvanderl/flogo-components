---
title: 3DS Get Classification Libraries
---

# 3DS Login
This activity allows you to get Classification Libraries from the Dassault Systemes 3DEXPERIENCE Platform.

## Installation
### Flogo Web
https://github.com/jvanderl/flogo-components/activity/dsgetlibraries
### Flogo CLI
```bash
flogo add activity github.com/jvanderl/flogo-components/activity/dsgetlibraries
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "3DServiceURL",
      "type": "string",
      "required": true
    },
    {
      "name": "accessToken",
      "type": "string",
      "required": true
    },
    {
      "name": "skipSsl",
      "type": "boolean",
      "value": false
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "any"
    },
    {
      "name": "status",
      "type": "integer"
    }
  ]
}
```
## Inputs
| Input     | Required | Description |
|:------------|:---------|:------------|
| 3DServiceURL   | True     | The URL for the service area to call |
| accessToken  | True    | The service access token (get from 3dslogin) |
| skipSsl     | False    | If set to true, skips the SSL validation (defaults to false)

## Outputs
| Output | Description |
|:------------|:------------|
| result  | The service response |         
| status       | Login Status, 0 if all is OK |

## Examples
### Simple
The below example logs in to 3DEXPERIENCE platform as user 'demoleader' for service 'flogo' to access services from '3DSpace':

```json
{
  "id": "dsgetlibraries",
  "name": "Get Libraries",
  "description": "Get Classification Libraroes from 3DS Platform",
  "activity": {
    "ref": "github.com/jvanderl/flogo-components/activity/dsgetlibraries",
    "input": {
      "3DServiceURL": "http://localhost/3DSpace",
      "userName": "<token>",
      "skipSsl": "true"
    }
  }
}
```
