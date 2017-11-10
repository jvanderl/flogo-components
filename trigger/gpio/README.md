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
  ],
  "outputs": [
    {
      "name": "state",
      "type": "int"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "gpiopin",
        "type": "int",
        "required" : true
      },
      {
        "name": "state",
        "type": "int",
        "required" : false
      }
    ]
  }
}
```
## Settings
- None -

## Ouputs
| Output   | Description    |
|:---------|:---------------|
| state    | The timer parameters used to trigger this flow |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| gpiopin   | GPIO pin to monitor |
| state     | State to monitor, 0 for low, 1 for high |

## Example Configuration

Triggers are configured via the triggers.json of your application. The following is and example configuration of the GPIO Trigger.

### Only once and immediate
Configure the Trigger to run a flow when pin 7 becomes high
```json
{
  "name": "gpio",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "gpiopin": "7",
				"state": "1"
      }
    }
  ]
}
```
