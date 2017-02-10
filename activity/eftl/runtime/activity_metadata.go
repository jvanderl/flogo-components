package eftl

var jsonMetadata = `{
  "name": "eftl",
  "version": "0.0.1",
  "title": "Send eFTL Message",
  "description": "Sends a message to TIBCO eFTL",
  "inputs":[
    {
      "name": "Server",
      "type": "string",
      "value": ""
    },
    {
      "name": "Channel",
      "type": "string",
      "value": ""
    },
    {
      "name": "Destination",
      "type": "string",
      "value": ""
    },
    {
      "name": "Message",
      "type": "string",
      "value": ""
    },
    {
      "name": "Username",
      "type": "string",
      "value": ""
    },
    {
      "name": "Password",
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
