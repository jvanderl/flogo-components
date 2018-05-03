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

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("ClusterIP", "127.0.0.1")
	tc.SetInput("Keyspace", "sample")
	tc.SetInput("TableName", "employee")
	tc.SetInput("Select", "*")
	tc.SetInput("Where", "empid=103")

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
	result := tc.GetOutput("result")
	switch v := result.(type) {
	case []map[string]interface{}:
		for s, a := range v {
			fmt.Printf("%v: record=%v\n", s, a)
			assert.Equal(t, expected, a)
		}
	}

}
