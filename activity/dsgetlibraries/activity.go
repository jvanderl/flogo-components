package dsgetlibraries

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-jvanderl-dsgetlibraries")

var savedCookies []http.Cookie

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodPATCH  = "PATCH"
	methodDELETE = "DELETE"

	iv3DServiceURL = "3DServiceURL"
	ivAccessToken  = "accessToken"
	ivSkipSsl      = "skipSsl"

	ovResult = "result"
	ovStatus = "status"
)

// RESTActivity is an Activity that is used to invoke a REST Operation
// inputs : {method,uri,params}
// outputs: {result}
type RESTActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new RESTActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &RESTActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *RESTActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *RESTActivity) Eval(context activity.Context) (done bool, err error) {

	skipSsl, _ := context.GetInput(ivSkipSsl).(bool)
	serviceAccessToken, _ := context.GetInput(ivAccessToken).(string)

	//Call Service
	serviceResponse, err := callService(context.GetInput(iv3DServiceURL).(string),
		serviceAccessToken,
		skipSsl)
	if err != nil {
		return false, err
	}
	log.Infof("Service Response: %v", serviceResponse)

	context.SetOutput(ovResult, serviceResponse)
	context.SetOutput(ovStatus, "0")
	return true, nil
}

func callService(serviceURI string, serviceAccessToken string, skipSSL bool) (response interface{}, err error) {

	method := "GET"
	uri := serviceURI
	uri += "/resources/IPClassificationReuse/library/libraries"
	pathParams := make(map[string]string)
	pathParams["ticket"] = serviceAccessToken
	uri, err = BuildURI(uri, pathParams)
	if err != nil {
		return nil, err
	}
	log.Debugf("REST Call: [%s] %s\n", method, uri)
	var reqBody io.Reader
	contentType := "application/json; charset=UTF-8"
	reqBody = nil
	req, err := http.NewRequest(method, uri, reqBody)
	if err != nil {
		return nil, err
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	httpTransportSettings := &http.Transport{}
	var client *http.Client

	// Skip ssl validation
	if skipSSL {
		httpTransportSettings.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client = &http.Client{Transport: httpTransportSettings}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	log.Debugf("response Status: %v", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)

	var result interface{}

	d := json.NewDecoder(bytes.NewReader(respBody))
	d.UseNumber()
	err = d.Decode(&result)

	return result, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils
func saveCookies(resp *http.Response) {
	for i, cookie := range resp.Cookies() {
		log.Debugf("Saving Cookie [%v] '%v' : '%v'", i, cookie.Name, cookie.Value)
		savedCookies = append(savedCookies, *cookie)

	}
}

func restoreCookies(req *http.Request) {
	for i, cookie := range savedCookies {
		log.Debugf("Restoring Cookie [%v] '%v' : '%v'", i, cookie.Name, cookie.Value)
		req.AddCookie(&cookie)
	}
}

//todo just make contentType a setting
func getContentType(replyData interface{}) string {

	contentType := "application/json; charset=UTF-8"

	switch v := replyData.(type) {
	case string:
		if !strings.HasPrefix(v, "{") && !strings.HasPrefix(v, "[") {
			contentType = "text/plain; charset=UTF-8"
		}
	case int, int64, float64, bool, json.Number:
		contentType = "text/plain; charset=UTF-8"
	default:
		contentType = "application/json; charset=UTF-8"
	}

	return contentType
}

func stringInList(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}

// BuildURI is a URI builder
func BuildURI(uri string, values map[string]string) (builtURI string, err error) {

	var newURI *url.URL
	newURI, err = url.Parse(uri)
	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	for k, v := range values {
		parameters.Add(k, v)
	}
	newURI.RawQuery = parameters.Encode()

	log.Debugf("URI: '%v'", newURI.String())

	return newURI.String(), nil
}
