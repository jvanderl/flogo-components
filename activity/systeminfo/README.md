# Get System Info
This activity provides your flogo application the ability retreive Hostname and IP Address.


## Installation

```bash
flogo add activity github.com/jvanderl/flogo-components/activity/systeminfo
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[],
  "outputs": [
    {
      "name": "hostname",
      "type": "string"
    },
    {
      "name": "ipaddress",
      "type": "string"
    }
  ]
}
```
## Output Descriptions
| Output   | Description    |
|:----------|:---------------|
| hostname  | The Hostname of the machine running the flogo instance |
| ipaddress | The IP Address of the machine running the flogo instance |         


## Configuration Examples
### Simple
Configure a task in flow to retreive system information, then map the hostname returned to 'message' in a log task":

```json
{  
  "activityType":"systeminfo",
  "id":3,
  "name":"systeminfo",
  "type":1,
  "attributes":[]
},
{  
  "activityType":"tibco-log",
  "id":4,
  "name":"Logger",
  "type":1,
  "attributes":[  
     {  
        "name":"message",
        "value":"Message sent to eFTL Server",
        "type":"string"
     },
     {  
        "name":"flowInfo",
        "value":"true",
        "type":"boolean"
     },
     {  
        "name":"addToFlow",
        "value":"true",
        "type":"boolean"
     }
  ],
  "inputMappings": [
    { "type": 1, "value": "{A3.hostname}", "mapTo": "message" }
  ]         
}
```
