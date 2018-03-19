# jl-rest
This activity provides your flogo application the ability to invoke a REST service.


## Installation

```bash
flogo add activity github.com/jvanderl/flogo-components/incubator/activity/rest
```

## Schema
Inputs and Outputs:

```json
{
   "input":[
    {
      "name": "method",
      "type": "string",
      "required": true,
      "allowed" : ["GET", "POST", "PUT", "PATCH", "DELETE"]
    },
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "proxy",
      "type": "string",
      "required": false
    },
    {
      "name": "pathParams",
      "type": "params"
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "header",
      "type": "params"
    },
    {
      "name": "content",
      "type": "any"
    },
    {
      "name": "allowInsecure",
      "type": "boolean"
    },
    {
      "name": "useBasicAuth",
      "type": "boolean"
    },
    {
      "name": "userID",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
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
## Settings
| Setting     | Description    |
|:------------|:---------------|
| method      | The HTTP method to invoke |         
| uri         | The uri of the resource   |
| pathParams  | The path parameters |
| queryParams | The query parameters |
| header      | The header parameters |
| content     | The message content |
| allowInsecure | Skip security verification |
| useBasicAuth  | Enable basic authentication (fill in userID and password) |
| userID        | Basic authentication User ID |
| password      | Basic authentication Password |
Note: 

* **pathParams**: Is only required if you have params in your URI ( i.e. http://.../pet/:id )
* **content**: Is only used in POST, PUT, PATCH

## Configuration Examples
### Simple
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io):

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for pet 1234",
  "attributes": [
    { "name": "method", "value": "GET" },
    { "name": "uri", "value": "http://petstore.swagger.io/v2/pet/1234" }
  ]
}
```
### Using Path Params
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io) via parameters.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for Pet",
  "attributes": [
    { "name": "method", "value": "GET" },
    { "name": "uri", "value": "http://petstore.swagger.io/v2/pet/:id" },
    { "name": "params", "value": { "id": "1234"} }
  ]
}
```
### Advanced
Configure a task in flow to get pet from the [swagger petstore](http://petstore.swagger.io) using a flow attribute to specify the id.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for Pet",
  "attributes": [
    { "name": "method", "value": "GET" },
    { "name": "uri", "value": "http://petstore.swagger.io/v2/pet/:id" },
  ],
  "inputMappings": [
    { "type": 1, "value": "petId", "mapTo": "params.id" }
  ]
}
```
