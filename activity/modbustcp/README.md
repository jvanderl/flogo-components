# Modbus TCP Operation
This activity provides your flogo application the ability to interact with Modbus TCP servers.


## Installation

```bash
flogo install github.com/jvanderl/flogo-components/activity/modbustcp
```
Link for flogo web:
```
https://github.com/jvanderl/flogo-components/activity/modbustcp
```

## Schema
Inputs and Outputs:

```json
{
  "settings":[
    {
      "name": "server",
      "type": "string"
    },
    {
      "name": "timeout",
      "type": "integer"
    },
    {
      "name": "SlaveId",
      "type": "number"
    }
  ],
  "input":[

    {
      "name": "operation",
      "type": "string",
      "allowed" : ["ReadCoils", "ReadDiscreteInputs", "ReadInputRegisters", "ReadHoldingRegisters", "WriteSingleCoil", "WriteMultipleCoils", "WriteSingleRegister", "WriteMultipleRegisters"]
    },
    {
      "name": "address",
      "type": "integer"
    },
    {
      "name": "numElements",
      "type": "integer"
    },
    {
      "name": "data",
      "type": "any"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "data",
      "type": "any"
    }
  ]
}

```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| server    | The Modbus Server Address and Port ([hostname]:[port])
| timeout   | Timeout when connecting to server (seconds) |         
| slaveId   | The ID for that identifies your flogo app as Modbus Slave |

## Inputs
| Input   | Description    |
|:--------|:---------------|
| operation   | Operation to perform ("ReadCoils", "ReadDiscreteInputs", "ReadInputRegisters", "ReadHoldingRegisters", "WriteSingleCoil", "WriteMultipleCoils", "WriteSingleRegister", "WriteMultipleRegisters") |
| address     | The address where to perform the operation |
| numElements | The number of coils/inputs/registers to be read/written |
| data        | The data needed to perform the operation (only used for write operations) |

## Outputs
| Output  | Description    |
|:--------|:---------------|
| result  | The result status of the operation sending. "OK" when success |
| data    | The data returned from the operation  |

## Configuration Examples
### Simple
Configure a task in flow to read value from 2 Coils on address 300 from Modbus server at 127.0.0.1 at port 502:

```json
{
  "id": "modbus_3",
  "name": "Modbus TCP Operation",
  "description": "Performs operations on a Modbus TCP server",
  "activity": {
    "ref": "#modbustcp",
    "input": {
      "operation": "ReadCoils",
      "numElements": 2,
      "address": "300"
    },
    "settings": {
      "server": "127.0.0.1:502",
      "timeout": 10,
      "slaveId": 0xFF
    }
  }
}
```
