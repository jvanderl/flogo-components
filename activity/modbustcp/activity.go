package modbustcp

import (
	"encoding/binary"
	"math"
	"strconv"

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
	ctx.Logger().Debugf("Settings SlaveID: %s", settings.SlaveID)

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

	err = act.conn.handler.Connect()
	if err != nil {
		ctx.Logger().Errorf("Error reconnecting to Modbus Server: %v\n", err)
		return false, err
	}

	output := &Output{}
	var results interface{}
	switch input.Operation {
	case "ReadCoils":
		ctx.Logger().Infof("Reading %v Coil(s) from address %v", input.NumElements, input.Address)
		results, err = act.conn.client.ReadCoils(input.Address, input.NumElements)
		break
	case "ReadDiscreteInputs":
		ctx.Logger().Infof("Reading %v Discete Input(s) from address %v", input.NumElements, input.Address)
		results, err = act.conn.client.ReadDiscreteInputs(input.Address, input.NumElements)
		break
	case "ReadInputRegisters":
		ctx.Logger().Infof("Reading %v Input Register(s) from address %v", input.NumElements, input.Address)
		results, err = act.conn.client.ReadInputRegisters(input.Address, input.NumElements)
		break
	case "ReadHoldingRegisters":
		ctx.Logger().Infof("Reading %v Holding Register(s) from address %v", input.NumElements, input.Address)
		results, err = act.conn.client.ReadHoldingRegisters(input.Address, input.NumElements)
		break
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
		if coilValue == 1 {
			coilValue = 0xFF00
		}
		results, err = act.conn.client.WriteSingleCoil(input.Address, coilValue)
		break
	case "WriteMultipleCoils":
		ctx.Logger().Infof("Writing %v to %v Coil(s) at address %v", input.Data, input.NumElements, input.Address)
		coilValues, ok := input.Data.([]byte)
		if !ok {
			ctx.Logger().Debugf("Error converting input.Data to byte array\n")
			output = &Output{Result: "INVALID INPUT DATA"}
			err = ctx.SetOutputObject(output)
			return false, nil
		}
		results, err = act.conn.client.WriteMultipleCoils(input.Address, input.NumElements, coilValues)
		break
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
		results, err = act.conn.client.WriteSingleRegister(input.Address, registerValue)
		break
	case "WriteMultipleRegisters":
		ctx.Logger().Infof("Writing %v to %v Register(s) at address %v", input.Data, input.NumElements, input.Address)
		registerValues, ok := input.Data.([]byte)
		if !ok {
			ctx.Logger().Debugf("Error converting input.Data to byte array\n")
			output = &Output{Result: "INVALID INPUT DATA"}
			err = ctx.SetOutputObject(output)
			return false, nil
		}
		results, err = act.conn.client.WriteMultipleRegisters(input.Address, input.NumElements, registerValues)
		break
	default:
		ctx.Logger().Debugf("Unknown Operation: %s", input.Operation)
		output = &Output{Result: "UNKNOWN OPERATION"}
		err = ctx.SetOutputObject(output)
		return false, nil
	}
	if err != nil {
		ctx.Logger().Debugf("Error performing Modbus Operation %v: %v\n", input.Operation, err)
		output = &Output{Result: "ERROR"}
		err = ctx.SetOutputObject(output)
		return false, err
	} else {
		outData := convertResult(results, input.ReturnAs, input.NumElements)
		ctx.Logger().Infof("Return Data: %v\n", outData)
		output = &Output{Result: "OK", Data: outData}
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func convertResult(result interface{}, returnAs string, numElements uint16) interface{} {
	inData := result.([]byte)
	numBytes := len(inData)
	if (returnAs == "Int16" || returnAs == "Uint16") && numBytes%2 != 0 {
		return "Datatype " + returnAs + " requires sets of 2 bytes response data, got " + strconv.Itoa(numBytes)
	}
	if (returnAs == "Int32" || returnAs == "Uint32" || returnAs == "Float32") && numBytes%4 != 0 {
		return "Datatype " + returnAs + " requires sets of 4 bytes response data, got " + strconv.Itoa(numBytes)
	}
	if (returnAs == "Int64" || returnAs == "Uint64" || returnAs == "Float64") && numBytes%8 != 0 {
		return "Datatype " + returnAs + " requires sets of 8 bytes response data, got " + strconv.Itoa(numBytes)
	}
	switch returnAs {
	case "Default", "":
		return result
	case "Bool":
		return byteArrayToBoolArray(inData, numElements)
	case "Int8":
		return byteArrayToInt8(inData)
	case "Int16":
		return byteArrayToInt16(inData)
	case "Int32":
		return byteArrayToInt32(inData)
	case "Int64":
		return byteArrayToInt64(inData)
	case "Uint16":
		return byteArrayToUint16(inData)
	case "Uint32":
		return byteArrayToUint32(inData)
	case "Uint64":
		return byteArrayToUint64(inData)
	case "Float32":
		return byteArrayToFloat32(inData)
	case "Float64":
		return byteArrayToFloat64(inData)
	default:
		return "Unknown Data type: " + returnAs
	}

}

func byteArrayToBoolArray(inData []byte, numBools uint16) []bool {
	result := []bool{}
	for i := len(inData) - 1; i > -1; i-- {
		for j := 1; j < 256; j = j * 2 {
			result = append(result, inData[i]&byte(j) != 0)
			if uint16(len(result)) == numBools {
				break
			}
		}
	}
	return result
}

func byteArrayToInt8(inData []byte) []int8 {
	result := []int8{}
	for i := 0; i < len(inData); i++ {
		result = append(result, int8(inData[i]))
	}
	return result
}

func byteArrayToUint16(inData []byte) []uint16 {
	result := []uint16{}
	for i := 0; i < len(inData); i = i + 2 {
		result = append(result, binary.BigEndian.Uint16(inData[i:]))
	}
	return result
}

func byteArrayToInt16(inData []byte) []int16 {
	result := []int16{}
	for i := 0; i < len(inData); i = i + 2 {
		result = append(result, int16(binary.BigEndian.Uint16(inData[i:])))
	}
	return result
}

func byteArrayToUint32(inData []byte) []uint32 {
	result := []uint32{}
	for i := 0; i < len(inData); i = i + 4 {
		result = append(result, binary.BigEndian.Uint32(inData[i:]))
	}
	return result
}

func byteArrayToInt32(inData []byte) []int32 {
	result := []int32{}
	for i := 0; i < len(inData); i = i + 4 {
		result = append(result, int32(binary.BigEndian.Uint32(inData[i:])))
	}
	return result
}

func byteArrayToUint64(inData []byte) []uint64 {
	result := []uint64{}
	for i := 0; i < len(inData); i = i + 8 {
		result = append(result, binary.BigEndian.Uint64(inData[i:]))
	}
	return result
}

func byteArrayToInt64(inData []byte) []int64 {
	result := []int64{}
	for i := 0; i < len(inData); i = i + 8 {
		result = append(result, int64(binary.BigEndian.Uint64(inData[i:])))
	}
	return result
}

func byteArrayToFloat32(inData []byte) []float32 {
	result := []float32{}
	for i := 0; i < len(inData); i = i + 4 {
		result = append(result, math.Float32frombits(binary.BigEndian.Uint32(inData[i:])))
	}
	return result
}

func byteArrayToFloat64(inData []byte) []float64 {
	result := []float64{}
	for i := 0; i < len(inData); i = i + 8 {
		result = append(result, math.Float64frombits(binary.BigEndian.Uint64(inData[i:])))
	}
	return result
}
