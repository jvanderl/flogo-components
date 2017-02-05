package filter

var jsonMetadata = `{
  "name": "filter",
  "version": "0.0.1",
  "description": "Filters input data based on min, max value",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "inputs":[
    {
      "name": "input",
      "type": "string",
      "required": true
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
      "name": "inverse",
      "type": "boolean"
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
    }
  ]
}`