# Timert
This trigger provides your flogo application the ability to schedule a flow via scheduling service
It's a modified version of the original timer activity on https://github.com/TIBCOSoftware/flogo-contrib

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/incubator/trigger/timert
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/incubator/trigger/timert
```

## Schema
Outputs and Endpoint:

```json
{
  "settings": [
  ],
  "output": [
    {
      "name": "triggerTime",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "repeating",
        "type": "boolean",
        "required" : true
      },
      {
        "name": "hours",
        "type": "number",
        "required" : false
      },
      {
        "name": "minutes",
        "type": "number",
        "required" : false
      },
      {
        "name": "seconds",
        "type": "number",
        "required" : false
      },
      {
        "name": "startImmediate",
        "type": "boolean",
        "required" : false
      },
      {
        "name": "startDate",
        "type": "string",
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
|:----------|:---------------|
| triggerTime |  The date and time the trigger fired |

## Handlers
| Setting   | Description    |
|:----------|:---------------|
| repeating    | Make trigger fire flows repeatedly. If set to true, fill in any combination of hours, minutes and seconds below |
| hours    | Repeating interval in hours (can be combined with minutes and seconds)|
| minutes    | Repeating interval in minutes (can be combined with hours and seconds) |
| seconds    | Repeating interval in seconds (can be combined with hours and minutes)|
| startImmediate | When set to true, the trigger will fire directly, otherwise wait for the startDate stated below |
| startDate    | Start date for first trigger in RFC3339 format ("2006-01-02T15:04:05Z07:00") |

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Timer Trigger.

### Only once and immediate
Configure the Trigger to run a flow once and start immediately

```json
{
  "name": "timer",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "false",
				"startImmediate": "true"
      }
    }
  ]
}
```

### Only once at schedule time
Configure the Trigger to run a flow at a certain date/time. "startDate" settings format = "yyyy-mm-ddThh:MM:ssZhh:mm"

```json
{
  "name": "timer",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "false",
        "startImmediate": "false",
				"startDate" : "2017-06-14T01:41:00Z02:00"
      }
    }
  ]
}
```

### Repeating
Configure the Trigger to run a flow repeating every hours|minutes|seconds. If "startImmediate" set to false, the trigger will not fire immediately.  In this case the first execution will occur in 24 hours. If set to true the first execution will will occur immediately.

```json
{
  "name": "timer",
  "settings": {
  },
  "handlers": [
  	{
  		"actionId": "local://testFlow",
  		"settings": {
  			"repeating": "true",
        "hours": "24",
        "minutes": "0",
  			"seconds": "0",
        "startImmediate": "true"
  		}
  	}
  ]
}
```

### Repeating with start date
Configure the Trigger to run a flow at a certain date/time and from then repeating every hours|minutes|seconds

```json
{
  "name": "timer",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "true",
        "hours": "1",
        "minutes": "10",
				"seconds": "5",
        "startImmediate": "false",
        "startDate" : "2017-06-14T2:28:00Z02:00"
      }
    }
  ]
}
```
