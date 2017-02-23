package eftl

var jsonMetadata = `{
  "name": "eftl",
  "version": "0.0.1",
  "description": "eFTL Subscriber",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "settings":[
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
      "name": "user",
      "type": "string",
      "value": ""
    },
    {
      "name": "password",
      "type": "string",
      "value": ""
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "destination",
        "type": "string"
      }
    ]
  }
}`