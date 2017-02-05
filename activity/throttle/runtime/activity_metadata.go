package throttle

var jsonMetadata = `{
  "name": "throttle",
  "version": "0.0.1",
  "description": "Trottles output based on unique data source and time interval",
  "title": "Throttle Output",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "homepage": "https://github.com/jvanderl/flogo-components/activity/throttle",
  "inputs":[
    {
      "name": "datasource",
      "type": "string"
    },
    {
      "name": "interval",
      "type": "integer"
    },
    {
      "name": "intervaltype",
      "type": "string",
      "allowed" : ["hours", "minutes", "seconds", "milliseconds"]      
    }    
  ],
  "outputs": [
    {
      "name": "pass",
      "type": "boolean"
    },
    {
      "name": "reason",
      "type": "string"
    },
    {
      "name": "lasttimepassed",
      "type": "string"
    }
  ]
}`