# timer2
This trigger provides your flogo application the ability to schedule a flow via scheduling service
It's a modified version of the original timer activity on https://github.com/TIBCOSoftware/flogo-contrib

## Installation

```bash
flogo install github.com/jvanderl/flogo-components/trigger/timer2
```
Link for flogo web: https://github.com/jvanderl/flogo-components/trigger/timer2
```

## Schema
Outputs and Endpoint:

```json
{
  "settings": [
  ],
  "outputs": [
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
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
        "name": "startImmediate",
        "type": "boolean",
        "required" : false
      },
      {
        "name": "startDate",
        "type": "string",
        "required" : false
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
      }
    ]
  }
}
```

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Timer Trigger.

### Only once and immediate
Configure the Trigger to run a flow once and start immediately

```json
{
  "name": "timer2",
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
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "false",
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
  "name": "timer2",
  "settings": {
  },
  "handlers": [
  	{
  		"actionId": "local://testFlow",
  		"settings": {
  			"startImmediate": "true",
  			"repeating": "true",
  			"seconds": "0",
  			"minutes": "0",
  			"hours": "24"
  		}
  	}
  ]
}
```

### Repeating with start date
Configure the Trigger to run a flow at a certain date/time and repeating every hours|minutes|seconds

```json
{
  "name": "timer2",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "true",
        "startImmediate": "true",
        "startDate" : "2017-06-14T2:28:00Z02:00",
				"seconds": "5",
				"minutes": "10",
        "hours": "1"
      }
    }
  ]
}
```
