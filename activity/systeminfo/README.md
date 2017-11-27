# Get System Info
This activity provides your flogo application the ability retreive Hostname and IP Address.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/systeminfo
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/systeminfo
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "includenetmask",
      "type": "boolean",
      "value": "false"
    }
  ],
  "outputs": [
    {
      "name": "hostname",
      "type": "string"
    },
    {
      "name": "ipaddress",
      "type": "string"
    },
    {
      "name": "ip6address",
      "type": "string"
    },
    {
      "name": "macaddress",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| includenetmask      | include the netmask in ip addresses (/xx suffix) |

## Output Descriptions
| Output   | Description    |
|:----------|:---------------|
| hostname  | The Hostname of the machine running the flogo instance |
| ipaddress | The IP Address of the machine running the flogo instance |         
| ip6address | The IPv6 Address of the machine running the flogo instance |         
| macaddress | The MAC Address of the machine running the flogo instance |         


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
