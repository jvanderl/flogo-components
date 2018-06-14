package udp

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

// Run Once, Start Immediately
const testConfig string = `{
  "name": "udp",
  "settings": {
		"port": 22600,
		"multicast_group": "224.192.32.19"
  },
  "handlers": [
    {
      "actionId": "nextaction",
      "settings": {
        "handler_setting": "xxx"
      }
    }
  ]
}`

// Listen for F1-2017 data
const testConfig2 string = `{
  "name": "udp",
  "settings": {
		"port": 20777,
		"multicast_group": ""
  },
  "handlers": [
    {
      "actionId": "NextAction",
      "settings": {
        "handler_setting": "xxx"
      }
    }
  ]
}`

//192.168.1.19
type TestRunner struct {
}

type initContext struct {
	handlers []*trigger.Handler
}

//var Test action.Runner

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Infof("Ran Action (Run): %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Infof("Ran Action (RunAction): %v", act)
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	log.Infof("Ran Action (Execute): %v", act)
	return nil, nil
}
func TestInit(t *testing.T) {
	// New factory
	md := trigger.NewMetadata(getJsonMetadata())
	f := NewFactory(md)

	// New Trigger
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr := f.New(&config)
	runner := &TestRunner{}
	//initCtx := &initContext{handlers: make([]*trigger.Handler, 0, len(config.Handlers))}

	initCtx := &struct {
				handlers []*trigger.Handler
			}{
				handlers: make([]*trigger.Handler, 0, len(config.Handlers))
			}

	newTrg, isNew := tgr.(trigger.Initializable)
	newTrg.Initialize(initCtx)
}

/*
//TODO: Fix Test
func TestUDPTrigger(t *testing.T) {
	log.Info("Testing UDP")
	config := trigger.Config{}
	runner := &TestRunner{}
	json.Unmarshal([]byte(testConfig2), &config)
	triggerFactory := &udpTriggerFactory{}
	triggerFactory.metadata = trigger.NewMetadata(jsonMetadata)
	trg := triggerFactory.New(&config)
	if trg == nil {
		log.Errorf("cannot create Trigger nil for id '%v'", &config.Id)
	}

	log.Infof("Number of Handlers: %v", len(config.Handlers))
	log.Infof("Hander 1 action ref : %v", config.Handlers[0].ActionId)
	//initCtx := trigger.InitContext{handlers: config.Handlers}
	//test := make([]*trigger.Handler, 0, len(config.Handlers))
	//log.Infof("test: %v", test)
	//var InitCtx trigger.InitContext = &trigger.InitContext{handlers: test}
	//initCtx := &trigger.InitContext{handlers: test}
	initCtx := &initContext{handlers: make([]*trigger.Handler, 0, len(config.Handlers))}

	newTrg, isNew := trg.(trigger.Initializable)
	log.Infof("newTrg: %v", newTrg)
	log.Infof("isNew: %v", isNew)

	//create handlers for that trigger and init
	for _, hConfig := range config.Handlers {

		/*
			//create the action
			actionFactory := action.GetFactory(hConfig.Action.Ref)
			if actionFactory == nil {
				log.Errorf("Action Factory '%s' not registered", hConfig.Action.Ref)
			}

			act, err := actionFactory.New(hConfig.Action)
			if err != nil {
				log.Errorf("Error creating actionFactory: %v", err)
			}

		log.Infof("hConfig: %v", hConfig)
		log.Infof("trg.Metadata().Output: %v", trg.Metadata().Output)
		log.Infof("trg.Metadata().Reply: %v", trg.Metadata().Reply)
		log.Infof("runner: %v", runner)

		handler := &trigger.Handler{config: &hconfig, act: nil, outputMd: nil, replyMd: nil, runner: runner}
		//handler := trigger.NewHandler(hConfig, nil, trg.Metadata().Output, trg.Metadata().Reply, runner)
		//		handler := trigger.NewHandler(hConfig, act, trg.Metadata().Output, trg.Metadata().Reply, runner)
		initCtx.handlers = append(initCtx.handlers, handler)

	}

	//	newTrg.Initialize(initCtx)
	trg.Start()
	defer trg.Stop()
	log.Infof("Press CTRL-C to quit")
	for {
	}
}

*/
