# UDP
This trigger provides your flogo application a stream of UDP data from the specificed port

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/incubator/trigger/udp
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/incubator/trigger/udp
```

## Schema
Outputs and Endpoint:

```json
{
"settings":[
    {
      "name": "port",
      "type": "integer"
    },
    {
      "name": "multicast_group",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "payload",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "handler_setting",
        "type": "string"
      }
    ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| port      | port to listen on |
| multicast_group    | listen group for Mukticast messages |

## Ouputs
| Output   | Description    |
|:---------|:---------------|
| payload    | The raw data from the message |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| N/A       | awaiting better understanding  |


## Example Configuration

Triggers are configured via the triggers.json of your application. The following is and example configuration of the UDP Trigger.

### Read UDP Data 
Configure the Trigger to capture all data on a given port 
```json
{
  "name": "udp",
  "settings": {
		"port": 20777,
		"multicast_group": ""
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "handler_setting": "xxx"
      }
    }
  ]
}}
```
