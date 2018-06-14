package cassandra

/*
to make the test succed, make sure you have something similar setup in cassandra:

CREATE KEYSPACE sample WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

USE sample;

CREATE TABLE employee (
	EmpID int,
	Name text,
	Salary double,
	PRIMARY KEY(EmpID)
);


INSERT INTO employee (empID, Name, Salary)
  VALUES (103, 'pqr', 7000.50);
*/

import (
	"fmt"
	"io/ioutil"
	"testing"
	//"strconv"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestOK(t *testing.T) {

	fmt.Println("Setting up to succeed.")
	fmt.Println("Expecting rowCount=1 and result='[map[empid:103 name:pqr salary:7000.5]]'")

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("clusterIP", "127.0.0.1")
	tc.SetInput("keySpace", "sample")
	tc.SetInput("tableName", "employee")
	tc.SetInput("select", "*")
	tc.SetInput("where", "empid=103")

	act.Eval(tc)

	//check result attr

	var (
		empid  = 103
		name   = "pqr"
		salary = 7000.5
	)

	expected := make(map[string]interface{})
	expected["empid"] = empid
	expected["name"] = name
	expected["salary"] = salary
	rowCount := tc.GetOutput("rowCount")
	fmt.Printf("rowCount: %v\n", rowCount)
	result := tc.GetOutput("result")
	fmt.Printf("result: %v\n", result)

	assert.Equal(t, 1, rowCount)

	switch v := result.(type) {
	case []map[string]interface{}:
		for s, a := range v {
			fmt.Printf("%v: record=%v\n", s, a)
			assert.Equal(t, expected, a)
		}
	}

}

func TestConnectError(t *testing.T) {

	fmt.Println("Setting up to fail connection.")
	fmt.Println("Expecting rowCount=0 and result='ERROR_CONNECT'")

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("clusterIP", "127.0.0.2")
	tc.SetInput("keySpace", "sample")
	tc.SetInput("tableName", "employee")
	tc.SetInput("select", "*")
	tc.SetInput("where", "empid=103")

	act.Eval(tc)

	//check result attr

	rowCount := tc.GetOutput("rowCount")
	fmt.Printf("rowCount: %v\n", rowCount)
	result := tc.GetOutput("result")
	fmt.Printf("result: %v\n", result)

	assert.Equal(t, 0, rowCount)
	assert.Equal(t, "ERROR_CONNECT", result)

}

func TestNoData(t *testing.T) {

	fmt.Println("Setting up to return no data.")
	fmt.Println("Expecting rowCount=0 and result='NO_DATA'")

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("clusterIP", "127.0.0.1")
	tc.SetInput("keySpace", "sample")
	tc.SetInput("tableName", "employeeeeee")
	tc.SetInput("select", "*")
	tc.SetInput("where", "empid=103")

	act.Eval(tc)

	//check result attr

	rowCount := tc.GetOutput("rowCount")
	fmt.Printf("rowCount: %v\n", rowCount)
	result := tc.GetOutput("result")
	fmt.Printf("result: %v\n", result)

	assert.Equal(t, 0, rowCount)
	assert.Equal(t, "NO_DATA", result)

}
