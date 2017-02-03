package mqtt2

var jsonMetadata = `{
  "name": "mqtt2",
  "version": "0.0.1",
  "description": "MQTT Topic Subscriber",
  "settings":[
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
      "name": "store",
      "type": "string"
    },
    {
      "name": "topic",
      "type": "string"
    },
    {
      "name": "qos",
      "type": "number"
    },
    {
      "name": "cleansess",
      "type": "boolean"
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "actualtopic",
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
