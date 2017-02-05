package filter

var jsonMetadata = `{
  "name": "filter",
  "version": "0.0.1",
  "description": "Filters input data based on min, max value and time interval",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "inputs":[
    {
      "name": "input",
      "type": "string",
      "required": true
    },
    {
      "name": "datasource",
      "type": "string"
    },
    {
      "name": "datatype",
      "type": "string",
      "required": true,
      "allowed" : ["int", "uint", "float32"]
    },
    {
      "name": "minvalue",
      "type": "string"
    },
    {
      "name": "maxvalue",
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
      "name": "result",
      "type": "string"
    },
    {
      "name": "usevalue",
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