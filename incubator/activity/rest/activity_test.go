package rest

import (
	"encoding/json"
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

const reqPostStr string = `{
  "name": "my pet"
}
`

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

var petID string

func TestSimplePost(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("method", "POST")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet")
	tc.SetInput("content", reqPostStr)

	//eval
	act.Eval(tc)
	val := tc.GetOutput("result")

	fmt.Printf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	fmt.Println("petID:", petID)
}

func TestSimpleGet(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/1")

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestSimpleGetWithHeaders(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/1")

	headers := make(map[string]string)
	headers["TestHeader"] = "TestValue"
	tc.SetInput("header", headers)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

/*
// TODO fix this test

func TestParamGet(t *testing.T) {

	act := activity.Get("github.com/TIBCOSoftware/flogo-contrib/activity/rest")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/:id")

	pathParams := map[string]string{
		"id": petID,
	}
	tc.SetInput("pathParams", pathParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}
*/

func TestSimpleGetQP(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/findByStatus")

	queryParams := map[string]string{
		"status": "ava",
	}
	tc.SetInput("queryParams", queryParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestBuildURI(t *testing.T) {

	uri := "http://localhost:7070/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI2(t *testing.T) {

	uri := "https://127.0.0.1:7070/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id":  "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI3(t *testing.T) {

	uri := "http://localhost/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI4(t *testing.T) {

	uri := "https://127.0.0.1/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id":  "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestSimpleGetWithProxy(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("proxy", "http://localhost:12345")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/16")

	//eval
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestTDVBD(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	uri := "https://ec2-34-243-92-189.eu-west-1.compute.amazonaws.com:9512/rest/v2/propertyGroups/:group/properties/:property/resources"

	pathParams := map[string]string{
		"group":    "GroupJL",
		"property": "PropertyJL",
	}

	newURI := BuildURI(uri, pathParams)

	fmt.Println(newURI)

	var content = `{
		"resourceBeans": [
		  {
			"resourcePath": "/localhost_9410/services/databases/IntroDemo",
			"resourceType": "database"
		  }
		],
		"propertyValue": [
		  "Muhahahah from FLOGO!"
		]
	  }`

	tc.SetInput("method", "POST")
	tc.SetInput("uri", uri)
	tc.SetInput("pathParams", pathParams)
	tc.SetInput("content", content)
	tc.SetInput("allowInsecure", "true")
	tc.SetInput("useBasicAuth", "true")
	tc.SetInput("userID", "admin")
	tc.SetInput("password", "admin")

	//eval
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)

}
