# GPIO
This trigger provides your flogo application the ability to start a flow upon a GPIO pin

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/gpio
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/trigger/gpio
```

## Schema
Outputs and Endpoint:

```json
{
  "settings": [
    {
      "name": "interval",
      "type": "int",
      "required" : true
    }    
  ],
  "outputs": [
    {
      "name": "state",
      "type": "number"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "gpiopin",
        "type": "number",
        "required" : true
      },
      {
        "name": "state",
        "type": "number",
        "required" : true
      },
      {
        "name": "pull",
        "type": "boolean",
        "required" : true
      }
    ]
  }
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| interval   | time between pin scans in Millisecond |

## Ouputs
| Output   | Description    |
|:---------|:---------------|
| state    | The atual reading of the pin (0 of 1) |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| gpiopin   | GPIO pin to monitor |
| state     | State to monitor, 0 for low, 1 for high |

## Example Configuration

Triggers are configured via the triggers.json of your application. The following is and example configuration of the GPIO Trigger.

### Only once and immediate
Configure the Trigger to run a flow when pin 7 becomes high, check every 0.5 seconds
```json
{
  "name": "gpio",
  "settings": {
    "interval": "500"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "gpiopin": "7",
				"state": "1",
        "pull": "true"
      }
    }
  ]
}
```
