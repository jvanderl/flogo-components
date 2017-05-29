# BLE Subscriber (Master)
This trigger provides your flogo application the ability to start a flow via Bluetooth Low Energy
Where Flogo is the BLE master, inititating connection to the BLE device.

## Installation

```bash
flogo add trigger github.com/jvanderl/flogo-components/trigger/blemaster
```
Link for flogo web: https://github.com/jvanderl/flogo-components/trigger/blemaster

## Schema
Settings, Outputs and Endpoint:

```json
{
  "name": "blemaster",
  "version": "0.0.1",
  "type": "flogo:trigger",
  "description": "Bluetooth Low Energy Master",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "settings":[
    {
      "name": "devicename",
      "type": "string"
    },
    {
      "name": "autoreconnect",
      "type": "boolean"
    },
    {
      "name": "reconnectinterval",
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
      "name": "notification",
      "type": "string"
    },
    {
      "name": "deviceid",
      "type": "string"
    },
    {
      "name": "localname",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "service",
        "type": "string"
      },
      {
        "name": "characteristic",
        "type": "string"
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| devicename    | The Local Name of the BLE device. Leave empty to discover just the services in the handler section |
| autoreconnect     | Automatically reconnect after device is disconnected   |
| reconnectinterval | How long to wait before attempting to reconnect |
| intervaltype      | Unit of time used by reconnectinterval |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| notification    | The data sent by the BLE device |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| service | The BLE Service that exposes the data |
| characteristic | The BLE characteristic under the service stated above |

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the BLE Master Trigger.

### Start a flow
Configure the Trigger to start "testFlow". So in this case the "handlers" "settings" "devicename" is emplty, so it will start "testFlow" flow when a message arrives on characteristic "ffe1" exposed under sevice "ffe0" from any device in this case.

```json
{
  "name": "blemaster",
  "settings": {
    "devicename": "",
		"autoreconnect": "true",
		"reconnectinterval": "5",
		"intervaltype": "seconds"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "service": "ffe0",
				"characteristic": "ffe1"
      }
    }
  ]
}
```
