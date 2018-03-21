# encoder
This activity provides your flogo application the ability to encode en decode strings using Base64, Base32.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/incubator/activity/encoder
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/incubator/activity/encoder
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "encoder",
      "type": "string",
      "required": true,
      "allowed" : ["BASE32", "BASE64", "HEX"]
    },
    {
      "name": "action",
      "type": "string",
      "required": true,
      "allowed" : ["ENCODE", "DECODE"]
    },
    {
      "name": "input",
      "type": "string",
      "required": true
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "status",
      "type": "string"
    }
  ]
}
```
## Inputs
| Input     | Description    |
|:------------|:---------------|
| encoder      | The Encoder type. Select BASE32, BASE64, HEX |         
| action         | Action to perform. Select ENCODE or DECODE   |
| input  | The string to be encoded/decoded |
## Outputs
| Output     | Description    |
|:------------|:---------------|
| result      | The encoded or decoded string |         
| status         | "OK" when all went fine, otherwise holds the error encountered.   |

