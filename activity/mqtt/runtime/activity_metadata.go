package mqtt

var jsonMetadata = `{
  "name": "mqtt",
  "version": "0.0.1",
  "title": "publish MQTT Message",
  "description": "Pubishes a message on an MQTT topic",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "inputs":[
   {
      "name": "broker",
      "type": "string"
    },
    {
      "name": "id",
      "type": "string"
    },
    {
      "name": "user",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
    },
    {
      "name": "topic",
      "type": "string"
    },
    {
      "name": "qos",
      "type": "integer",
      "allowed" : ["0", "1", "2"]
    },
    {
      "name": "message",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}`
