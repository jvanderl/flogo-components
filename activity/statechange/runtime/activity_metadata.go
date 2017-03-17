package statechange

var jsonMetadata = `{
  "name": "statechange",
  "version": "0.0.1",
  "description": "Detects state change for up to eight inputs",
  "title": "State Change",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "homepage": "https://github.com/jvanderl/flogo-components/activity/statechange",
  "inputs":[
    {
      "name": "datasource",
      "type": "string"
    },
    {
      "name": "input1",
      "type": "string"
    },
    {
      "name": "input2",
      "type": "string"
    },
    {
      "name": "input3",
      "type": "string"
    },
    {
      "name": "input4",
      "type": "string"
    },
    {
      "name": "input5",
      "type": "string"
    },
    {
      "name": "input6",
      "type": "string"
    },
    {
      "name": "input7",
      "type": "string"
    },
    {
      "name": "input8",
      "type": "string"
    }   
  ],
  "outputs": [
    {
      "name": "changed",
      "type": "boolean"
    },
    {
      "name": "flags",
      "type": "integer"
    },
    {
      "name": "result",
      "type": "string"
    }
  ]
}`
