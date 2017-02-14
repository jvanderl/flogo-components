package kafka

var jsonMetadata = `{
  "name": "kafka",
  "version": "0.0.1",
  "title": "Receive Kafka Message",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "description": "Kafka Topic Subscriber",
  "settings":[
    {
      "name": "server",
      "type": "string"
    },
    {
      "name": "configid",
      "type": "string"
    },
    {
      "name": "topic",
      "type": "string"
    },
    {
      "name": "partition",
      "type": "int"
    },
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
        "name": "topic",
        "type": "string"
      }
    ]
  }
}`
