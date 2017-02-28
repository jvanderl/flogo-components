package eftl

var jsonMetadata = `{
  "name": "eftl",
  "version": "0.0.1",
  "title": "Send eFTL Message",
  "description": "Sends a message to TIBCO eFTL",
  "inputs":[
    {
      "name": "server",
      "type": "string",
      "value": ""
    },
    {
      "name": "channel",
      "type": "string",
      "value": ""
    },
    {
      "name": "destination",
      "type": "string",
      "value": ""
    },
    {
      "name": "message",
      "type": "string",
      "value": ""
    },
    {
      "name": "user",
      "type": "string",
      "value": ""
    },
    {
      "name": "password",
      "type": "string",
      "value": ""
    },
    {
      "name": "secure",
      "type": "boolean",
      "value": "false"
    },
    {
      "name": "certificate",
      "type": "string",
      "value": ""
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}`
