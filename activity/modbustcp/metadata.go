package modbustcp

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Server  string `md:"server,required"`  // The Modbus Server Address and Port ([hostname]:[port])
	Timeout int    `md:"timeout,required"` // Timeout when connecting to server (seconds)
	SlaveID int    `md:"slaveId,required"` // The ID for that identifies your flogo app as Modbus Slave
}

type Input struct {
	Operation   string      `md:"operation,required"` // Operation to perform ("ReadCoils", "ReadDiscreteInputs", "ReadInputRegisters", "ReadHoldingRegisters", "WriteSingleCoil", "WriteMultipleCoils", "WriteSingleRegister", "WriteMultipleRegisters")
	Address     uint16      `md:"address,required"`   // The address where to perform the operation
	NumElements uint16      `md:"numElements"`        // The number of coils/inputs/registers to be read/written
	Data        interface{} `md:"data"`               // The data needed to perform the operation
	ReturnAs    string      `md:"returnAs"`           // Return Type ("Default", "Bool", "Int8", "Int16", "Int32", "Int64", "Uint16", "Uint32", "Uint64", "Float32", "Float64")
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"operation": i.Operation,
		"address":   i.Address,
		"data":      i.Data,
		"returnAs":  i.ReturnAs,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Operation, err = coerce.ToString(values["operation"])
	tmpInt, err := coerce.ToInt(values["address"])
	if err != nil {
		return err
	}
	i.Address = uint16(tmpInt)
	if val, ok := values["numElements"]; ok {
		tmpInt, err = coerce.ToInt(val)
		if err != nil {
			return err
		}
		i.NumElements = uint16(tmpInt)
	}
	if val, ok := values["data"]; ok {
		i.Data, err = coerce.ToAny(val)
		if err != nil {
			return err
		}
	}
	if val, ok := values["returnAs"]; ok {
		i.ReturnAs, err = coerce.ToString(val)
	} else {
		i.ReturnAs = "Default"
	}
	return nil
}

type Output struct {
	Result string      `md:"result"` // The result status of the operation sending. "OK" when success
	Data   interface{} `md:"data"`   // The data returned from the operation
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
		"data":   o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Result, _ = coerce.ToString(values["result"])
	o.Data, _ = coerce.ToAny(values["result"])

	return nil
}
