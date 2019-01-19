---
title: 3DS Login
---

# 3DS Login
This activity allows you to login to the Dassault Systemes 3DEXPERIENCE Platform.

## Installation
### Flogo Web
https://github.com/jvanderl/flogo-components/activity/dslogin
### Flogo CLI
```bash
flogo add activity github.com/jvanderl/flogo-components/activity/dslogin
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "3DPassportURL",
      "type": "string",
      "required": true
    },
    {
      "name": "3DServiceURL",
      "type": "string",
      "required": true
    },
    {
      "name": "userName",
      "type": "string",
      "required": true
    },
    {
      "name": "serviceName",
      "type": "string",
      "required": true
    },
    {
      "name": "serviceSecret",
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
      "name": "serviceAccessToken",
      "type": "string"
    },
    {
      "name": "serviceRedirectURL",
      "type": "string"
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
| 3DPassportURL | True     | the URL for the platform 3DPassport Service |         
| 3DServiceURL   | True     | The URL for the service area to call |
| userName       | True    | The User Name to logon the the platform |
| serviceName  | True    | The Service Name from 3DPassport Integration configuration|
| serviceSecret | True    | TThe Service Secret from 3DPassport Integration configuration |
| skipSsl     | False    | If set to true, skips the SSL validation (defaults to false)

## Outputs
| SettiOutput | Description |
|:------------|:---------|:------------|
| serviceAccessToken  | Use this token in subsequent REST calls to the platform |         
| serviceRedirectURL   | The URL to login to the service area with the token provided |
| status       | Login Status, 0 if all is OK |

## Examples
### Simple
The below example logs in to 3DEXPERIENCE platform as user 'demoleader' for service 'flogo' to access services from '3DSpace':

```json
{
  "id": "dslogin",
  "name": "Login to 3DS Platform",
  "description": "Login to 3DS Platform",
  "activity": {
    "ref": "github.com/jvanderl/flogo-components/activity/dslogin",
    "input": {
      "3DPassportURL": "http://localhost/3DPassport",
      "3DServiceURL": "http://localhost/3DSpace",
      "userName": "demoleader",
      "serviceName": "flogo",
      "serviceName": "<servicesecret>",
      "skipSsl": "true"
    }
  }
}
```
