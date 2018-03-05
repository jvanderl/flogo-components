# Check IBAN
This activity provides your flogo application the ability to check an International Bank Account Number.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/checkiban
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/checkiban
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "iban",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "printcode",
      "type": "string"
    },
    {
      "name": "code",
      "type": "string"
    },
    {
      "name": "countrycode",
      "type": "string"
    },
    {
      "name": "checkdigits",
      "type": "string"
    },
    {
      "name": "bban",
      "type": "string"
    },
    {
      "name": "ibanobj",
      "type": "object"
    }        
  ]
}

```
## Inputs
| Input     | Description    |
|:------------|:---------------|
| iban      | International Bank Account Number |

## Output Descriptions
| Output   | Description    |
|:----------|:---------------|
| result  | Validation result, OK when valid IBAN |
| printcode | Printable version of the IBAN code, when valid |         
| code |  IBAN code returned when valid |         
| countrycode  | The country code section of a valid IBAN |
| checkdigits | The check digits section of a valid IBAN |         
| bban |  The Basic Bank Account Number section of a valid IBAN |         
| ibanobj  | The complete returned object of a valid IBAN |
     

## Configuration Examples
### Simple
Configure a task in flow to check an IBAN number, then map the printcode returned to 'message' in a log task":

```json
{  
  "activityType":"checkiban",
  "id":3,
  "name":"checkiban",
  "type":1,
  "attributes":[
    {
        "name":"iban",
        "value":"NL40ABNA0517552264",
        "type":"string"
    }
  ]
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
    { "type": 1, "value": "{A3.printcode}", "mapTo": "message" }
  ]         
}
```
