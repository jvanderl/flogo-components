package modbustcp

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("Settings Server : %s", settings.Server)
	ctx.Logger().Debugf("Settings Timeout: %s", settings.Timeout)
	ctx.Logger().Debugf("Settings SlaveID	 : %s", settings.SlaveID)

	conn, err := getModbusTcpConnection(ctx.Logger(), settings)
	if err != nil {
		return nil, err
	}

	act := &Activity{conn: conn}

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	conn *ModbusTcpConnection
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	ctx.Logger().Debugf("Operation: %s", input.Operation)
	ctx.Logger().Debugf("Address: %v", input.Address)
	ctx.Logger().Debugf("NumElements: %v", input.NumElements)
	ctx.Logger().Debugf("Data: %v", input.Data)

	output := &Output{}

	switch input.Operation {
	case "ReadCoils":
		ctx.Logger().Infof("Reading %v Coil(s) from address %v", input.NumElements, input.Address)
		results, err := act.conn.client.ReadCoils(input.Address, input.NumElements)
		if err != nil {
			ctx.Logger().Debugf("Error reading coils from Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR READING COILS"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Coil Values: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "ReadDiscreteInputs":
		ctx.Logger().Infof("Reading %v Discete Input(s) from address %v", input.NumElements, input.Address)
		results, err := act.conn.client.ReadDiscreteInputs(input.Address, input.NumElements)
		if err != nil {
			ctx.Logger().Debugf("Error reading discrete inputs from Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR READING DISCRETE INPUTS"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Discrete Input Values: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "ReadInputRegisters":
		ctx.Logger().Infof("Reading %v Input Register(s) from address %v", input.NumElements, input.Address)
		results, err := act.conn.client.ReadInputRegisters(input.Address, input.NumElements)
		if err != nil {
			ctx.Logger().Debugf("Error reading input registers from Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR READING INPUT REGISTERS"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Input Register Values: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "ReadHoldingRegisters":
		ctx.Logger().Infof("Reading %v Holding Register(s) from address %v", input.NumElements, input.Address)
		results, err := act.conn.client.ReadHoldingRegisters(input.Address, input.NumElements)
		if err != nil {
			ctx.Logger().Debugf("Error reading holding registers from Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR READING HOLDING REGISTERS"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Holding Register Values: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "WriteSingleCoil":
		ctx.Logger().Infof("Writing %X to Single Coil at address %v", input.Data, input.Address)
		tmpInt, ok := input.Data.(int)
		if !ok {
			ctx.Logger().Debugf("Error converting input.Data to uint16\n")
			output = &Output{Result: "INVALID INPUT DATA"}
			err = ctx.SetOutputObject(output)
			return false, nil
		}
		coilValue := uint16(tmpInt)
		results, err := act.conn.client.WriteSingleCoil(input.Address, coilValue)
		if err != nil {
			ctx.Logger().Debugf("Error writing single coil to Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR WRITING SINGLE COIL"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Return Data: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "WriteMultipleCoils":
		ctx.Logger().Infof("Writing %v to %v Coil(s) at address %v", input.Data, input.NumElements, input.Address)
		coilValues, ok := input.Data.([]byte)
		if !ok {
			ctx.Logger().Debugf("Error converting input.Data to byte array\n")
			output = &Output{Result: "INVALID INPUT DATA"}
			err = ctx.SetOutputObject(output)
			return false, nil
		}
		results, err := act.conn.client.WriteMultipleCoils(input.Address, input.NumElements, coilValues)
		if err != nil {
			ctx.Logger().Debugf("Error writing multiple coils to Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR WRITING MULTIPLE COILS"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Return Data: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "WriteSingleRegister":
		ctx.Logger().Infof("Writing %X to Single Register at address %v", input.Data, input.Address)
		tmpInt, ok := input.Data.(int)
		if !ok {
			ctx.Logger().Debugf("Error converting input.Data to uint16\n")
			output = &Output{Result: "INVALID INPUT DATA"}
			err = ctx.SetOutputObject(output)
			return false, nil
		}
		registerValue := uint16(tmpInt)
		results, err := act.conn.client.WriteSingleRegister(input.Address, registerValue)
		if err != nil {
			ctx.Logger().Debugf("Error writing single register to Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR WRITING SINGLE REGISTER"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Return Data: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	case "WriteMultipleRegisters":
		ctx.Logger().Infof("Writing %v to %v Register(s) at address %v", input.Data, input.NumElements, input.Address)
		registerValues, ok := input.Data.([]byte)
		if !ok {
			ctx.Logger().Debugf("Error converting input.Data to byte array\n")
			output = &Output{Result: "INVALID INPUT DATA"}
			err = ctx.SetOutputObject(output)
			return false, nil
		}
		results, err := act.conn.client.WriteMultipleRegisters(input.Address, input.NumElements, registerValues)
		if err != nil {
			ctx.Logger().Debugf("Error writing multiple registers to Modbus Server: %v\n", err)
			output = &Output{Result: "ERROR WRITING MULTIPLE REGISTERS"}
			err = ctx.SetOutputObject(output)
			return false, nil
		} else {
			ctx.Logger().Infof("Return Data: %v\n", results)
			output = &Output{Result: "OK", Data: results}
		}
	default:
		ctx.Logger().Debugf("Unknown Operation: %s", input.Operation)
		output = &Output{Result: "UNKNOWN OPERATION"}
		err = ctx.SetOutputObject(output)
		return false, nil
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
