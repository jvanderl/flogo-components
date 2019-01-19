package dslogin

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-rest")

var savedCookies []http.Cookie

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodPATCH  = "PATCH"
	methodDELETE = "DELETE"

	iv3DPassportURL = "3DPassportURL"
	iv3DServiceURL  = "3DServiceURL"
	ivUserName      = "userName"
	ivServiceName   = "serviceName"
	ivserviceSecret = "serviceSecret"
	ivSkipSsl       = "skipSsl"

	ovServiceAccessToken = "serviceAccessToken"
	ovServiceRedirectURL = "serviceRedirectURL"
	ovStatus             = "status"
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

	//first get CAS login token
	log.Infof("Getting CAS Access Token for service '%v'", context.GetInput(ivServiceName).(string))

	accessToken, reAuthURL, err := getCASTransientToken(context.GetInput(iv3DPassportURL).(string),
		context.GetInput(ivUserName).(string),
		context.GetInput(ivServiceName).(string),
		context.GetInput(ivserviceSecret).(string),
		skipSsl)
	if err != nil {
		return false, err
	}

	log.Infof("CAS Access Token: %v", accessToken)
	log.Infof("ReAuth URL: %v", reAuthURL)

	// Then do CAS Login
	log.Infof("Performing CAS Login for URL '%v'", context.GetInput(iv3DServiceURL).(string))

	serviceAccessToken, serviceRedirectURL, err := doCASLoginForService(context.GetInput(iv3DPassportURL).(string),
		context.GetInput(iv3DServiceURL).(string),
		accessToken,
		skipSsl)
	if err != nil {
		return false, err
	}

	log.Infof("Service Access Token: %v", serviceAccessToken)
	log.Infof("Service Redirect URL: %v", serviceRedirectURL)

	//Now do Service Login
	/*loginResponse, err := doServiceLogin(context.GetInput(iv3DPassportURL).(string),
		context.GetInput(iv3DServiceURL).(string),
		serviceAccessToken, skipSsl)
	if err != nil {
		return false, err
	}
	log.Infof("Login Response: %v", loginResponse)
	*/
	//Now call Test Service
	/*serviceResponse, err := callTestService(context.GetInput(iv3DServiceURL).(string),
		serviceAccessToken,
		skipSsl)
	if err != nil {
		return false, err
	}
	log.Infof("Service Response: %v", serviceResponse)*/

	context.SetOutput(ovServiceAccessToken, serviceAccessToken)
	context.SetOutput(ovServiceRedirectURL, serviceRedirectURL)
	context.SetOutput(ovStatus, "0")
	return true, nil
}

func getCASTransientToken(passportURI string, userName string, serviceName string, serviceSecret string, skipSSL bool) (accessToken string, reAuthURL string, err error) {

	method := "GET"
	uri := passportURI
	uri += "/api/v2/batch/ticket"
	pathParams := make(map[string]string)
	pathParams["identifier"] = userName
	uri = BuildURI(uri, pathParams)
	log.Debugf("REST Call: [%s] %s\n", method, uri)
	var reqBody io.Reader
	contentType := "application/json; charset=UTF-8"
	reqBody = nil
	req, err := http.NewRequest(method, uri, reqBody)
	if err != nil {
		return "", "", err
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("DS-Service-Name", serviceName)
	req.Header.Set("DS-Service-Secret", serviceSecret)

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
		return "", "", err
	}
	log.Debugf("response Status: %v", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)

	saveCookies(resp)

	var result map[string]interface{}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", "", err
	}

	return result["access_token"].(string), result["x3ds_reauth_url"].(string), nil
}

func doCASLoginForService(passportURI string, serviceURI string, accessToken string, skipSSL bool) (svcAccessToken string, redirectURL string, err error) {

	method := "GET"
	uri := passportURI
	uri += "/api/login/cas/transient"
	pathParams := make(map[string]string)
	pathParams["tgt"] = accessToken
	pathParams["service"] = serviceURI
	uri = BuildURI(uri, pathParams)
	log.Debugf("REST Call: [%s] %s\n", method, uri)
	var reqBody io.Reader
	contentType := "application/json; charset=UTF-8"
	reqBody = nil
	req, err := http.NewRequest(method, uri, reqBody)
	if err != nil {
		return "", "", err
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	restoreCookies(req)

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
		return "", "", err
	}
	log.Debugf("response Status: %v", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)

	saveCookies(resp)

	var result map[string]interface{}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", "", err
	}

	return result["access_token"].(string), result["x3ds_service_redirect_url"].(string), nil
}

/*func doServiceLogin(passportURI string, serviceURI string, scvAccessToken string, skipSSL bool) (response interface{}, err error) {

	method := "GET"
	uri := passportURI
	uri += "/login"
	pathParams := make(map[string]string)
	pathParams["service"] = serviceURI
	pathParams["ticket"] = scvAccessToken
	uri = BuildURI(uri, pathParams)
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

	restoreCookies(req)

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
	log.Infof("response Status: %v", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)

	saveCookies(resp)

	var result interface{}

	d := json.NewDecoder(bytes.NewReader(respBody))
	d.UseNumber()
	err = d.Decode(&result)

	return result, nil
}*/

/*
func callTestService(serviceURI string, serviceAccessToken string, skipSSL bool) (response interface{}, err error) {

	method := "GET"
	uri := serviceURI
	uri += "/resources/IPClassificationReuse/library/libraries"
	pathParams := make(map[string]string)
	pathParams["ticket"] = serviceAccessToken
	uri = BuildURI(uri, pathParams)
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
*/
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

// BuildURI is a temporary crude URI builder
func BuildURI(uri string, values map[string]string) string {

	newURI := uri

	var paramIndex = 0
	for k, v := range values {
		paramIndex++
		if paramIndex == 1 {
			// first entry, use "?"
			newURI += "?"
		} else {
			newURI += "&"
		}
		newURI += k + "=" + v
	}

	log.Debugf("URI: '%v'", newURI)
	return newURI
}
