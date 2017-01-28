package flowprovider

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/script/fggos"
	"github.com/TIBCOSoftware/flogo-lib/flow/service"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("flowprovider")

const (
	uriSchemeFile     = "file://"
	uriSchemeEmbedded = "embedded://"
)

// RemoteFlowProvider is an implementation of FlowProvider service
// that can access flowes via URI
type RemoteFlowProvider struct {
	//todo: switch to LRU cache
	mutex       *sync.Mutex
	enabled     bool
	flowCache   map[string]*flowdef.Definition
	embeddedMgr *support.EmbeddedFlowManager
}

// NewRemoteFlowProvider creates a RemoteFlowProvider
func NewRemoteFlowProvider(config *util.ServiceConfig, embeddedFlowMgr *support.EmbeddedFlowManager) *RemoteFlowProvider {

	var service RemoteFlowProvider
	service.flowCache = make(map[string]*flowdef.Definition)
	service.mutex = &sync.Mutex{}
	service.enabled = config.Enabled
	service.embeddedMgr = embeddedFlowMgr
	return &service
}

func (pps *RemoteFlowProvider) Name() string {
	return service.ServiceFlowProvider
}

func (pps *RemoteFlowProvider) Enabled() bool {
	return pps.enabled
}

// Start implements util.Managed.Start()
func (pps *RemoteFlowProvider) Start() error {
	// no-op
	return nil
}

// Stop implements util.Managed.Stop()
func (pps *RemoteFlowProvider) Stop() error {
	// no-op
	return nil
}

// GetFlow implements flow.Provider.GetFlow
func (pps *RemoteFlowProvider) GetFlow(flowURI string) *flowdef.Definition {

	//handle panic

	// todo turn pps.flowCache to real cache
	if flow, ok := pps.flowCache[flowURI]; ok {
		log.Debugf("Accessing cached Flow: %s\n")
		return flow
	}

	log.Infof("Get Flow: %s\n", flowURI)

	var flowJSON []byte

	if strings.HasPrefix(flowURI, uriSchemeEmbedded) {

		log.Infof("Loading Embedded Flow: %s\n", flowURI)
		flowJSON = pps.embeddedMgr.GetEmbeddedFlowJSON(flowURI)

	} else if strings.HasPrefix(flowURI, uriSchemeFile) {

		log.Infof("Loading Local Flow: %s\n", flowURI)
		flowFilePath, _ := util.URLStringToFilePath(flowURI)

		flowJSON, _ = ioutil.ReadFile(flowFilePath)

	} else {

		log.Infof("Loading Remote Flow: %s\n", flowURI)

		req, err := http.NewRequest("GET", flowURI, nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err) //todo probably shouldn't panic
		}
		defer resp.Body.Close()

		log.Infof("response Status:", resp.Status)

		if resp.StatusCode >= 300 {
			//not found
			return nil
		}

		body, _ := ioutil.ReadAll(resp.Body)
		flowJSON = body
	}

	if flowJSON != nil {
		var defRep flowdef.DefinitionRep
		json.Unmarshal(flowJSON, &defRep)

		def, err := flowdef.NewDefinition(&defRep)

		if err != nil {
			log.Errorf("Error unmarshalling flow: %s", err.Error())
			log.Debugf("Failed to unmarshal: %s", string(flowJSON))

			return nil
		}

		//todo optimize this - not needed if flow doesn't have expressions
		//todo have a registry for this?
		def.SetLinkExprManager(fggos.NewGosLinkExprManager(def))
		//def.SetLinkExprManager(fglua.NewLuaLinkExprManager(def))

		//synchronize
		pps.mutex.Lock()
		pps.flowCache[flowURI] = def
		pps.mutex.Unlock()

		return def
	}

	log.Debugf("Flow Not Found: %s\n", flowURI)

	return nil
}

func DefaultConfig() *util.ServiceConfig {
	return &util.ServiceConfig{Name: service.ServiceFlowProvider, Enabled: true}
}
