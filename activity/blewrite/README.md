# Write BLE Data
This activity provides your Flogo application the ability to write data to a Bluetooth Low Energy (BLE) device.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/blewrite
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/blewrite
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
   {
      "name": "devicename",
      "type": "string"
    },
    {
      "name": "deviceid",
      "type": "string"
    },
    {
      "name": "serviceid",
      "type": "string"
    },
    {
      "name": "characteristic",
      "type": "string"
    },
    {
      "name": "bledata",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| devicename    | The BLE Local Name of the target device|
| deviceid        | The Device ID of the target device |         
| serviceid      | The BLE Service ID (group of characteristics) that you want to use on the target device |
| characteristic  | The BLE Charateristic you want to write the BLE data to |
| bledata     | The data you wish to send to the BLE device |



## Configuration Examples
### Simple
Configure a task in flow to write "hello world" to a BLE device named "IOTDEVICE" with deviceID "A4:D5:78:6D:57:6C":

```json
{
  "id": 2,
  "name": "Write BLE data",
  "type": 1,
  "activityType": "blewrite",
  "attributes": [
    {
      "name": "devicename",
      "value": "IOTDEVICE",
      "type": "string"
    },
    {
      "name": "deviceid",
      "value": "A4:D5:78:6D:57:6C",
      "type": "string"
    },
    {
      "name": "serviceid",
      "value": "ffe0",
      "type": "string"
    },
    {
      "name": "characteristic",
      "value": "ffe1",
      "type": "string"
    },
    {
      "name": "bledata",
      "value": "hello world",
      "type": "string"
    }
  ]
}
```
