package test

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("test")

func init() {
	model.Register(NewTestModel())
}

func NewTestModel() *model.FlowModel {
	m := model.New("test")
	m.RegisterFlowBehavior(&SimpleFlowBehavior{})
	m.RegisterTaskBehavior(1, &SimpleTaskBehavior{})

	return m
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleFlow

type SimpleFlowBehavior struct {
}

func (b *SimpleFlowBehavior) Start(context model.FlowContext) (start bool, evalCode int) {

	//just schedule the root task
	return true, 0
}

func (b *SimpleFlowBehavior) Resume(context model.FlowContext) bool {

	return true
}

func (b *SimpleFlowBehavior) TasksDone(context model.FlowContext, doneCode int) {
	log.Debugf("Flow TasksDone\n")

}

func (b *SimpleFlowBehavior) Done(context model.FlowContext) {
	log.Debugf("Flow Done\n")

}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleTask

type SimpleTaskBehavior struct {
}

func (b *SimpleTaskBehavior) Enter(context model.TaskContext, enterCode int) (eval bool, evalCode int) {

	task := context.Task()
	//check if all predecessor links are done
	log.Debugf("Task Enter: %s\n", task.Name())

	context.SetState(STATE_ENTERED)

	linkContexts := context.FromInstLinks()

	ready := true

	if len(linkContexts) == 0 {
		ready = true
	} else {

		log.Debugf("Num Links: %d\n", len(linkContexts))
		for _, linkContext := range linkContexts {

			log.Debugf("Task: %s, linkData: %v\n", task.Name(), linkContext)
			if linkContext.State() != STATE_LINK_TRUE {
				ready = false
				break
			}
		}
	}

	if ready {
		log.Debugf("Task Ready\n")
		context.SetState(STATE_READY)
	} else {
		log.Debugf("Task Not Ready\n")
	}

	return ready, 0
}

func (b *SimpleTaskBehavior) Eval(context model.TaskContext, evalCode int) (done bool, doneCode int, err error) {

	task := context.Task()
	log.Debugf("Task Eval: %s\n", task)

	if len(task.ChildTasks()) > 0 {
		log.Debugf("Has Children\n")

		context.SetState(STATE_WAITING)

		//for now enter all children (bpel style) - costly
		context.EnterChildren(nil)

		return false, 0, nil
	}

	if context.HasActivity() {

		//log.Debug("Evaluating Activity: ", activity.GetType())
		done, _ := context.EvalActivity()
		return done, 0, nil
	}

	//no-op
	return true, 0, nil
}

func (b *SimpleTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {
	log.Debugf("Task PostEval\n")

	if context.HasActivity() { //and is async

		//done := activity.PostEval(activityContext, data)
		done := true

		return done, 0, nil
	}
	//no-op
	return true, 0, nil
}

func (b *SimpleTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry) {

	context.SetState(STATE_DONE)
	//context.SetTaskDone() for task garbage collection

	task := context.Task()

	log.Debugf("done task:%s\n", task.Name())

	links := task.ToLinks()

	numLinks := len(links)

	if numLinks > 0 {

		taskEntries := make([]*model.TaskEntry, 0, numLinks)

		for _, link := range links {

			follow, _ := context.EvalLink(link)
			if follow {

				taskEntry := &model.TaskEntry{Task: link.ToTask(), EnterCode: 0}
				taskEntries = append(taskEntries, taskEntry)
			}
		}

		//continue on to successor links
		return false, 0, taskEntries
	}

	//notify parent that we are done
	return true, 0, nil
}

func (b *SimpleTaskBehavior) ChildDone(context model.TaskContext, childTask *flowdef.Task, childDoneCode int) (done bool, doneCode int) {
	log.Debugf("Task ChildDone\n")

	return true, 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleLink

type SimpleLinkBehavior struct {
}

func (b *SimpleLinkBehavior) Eval(context model.LinkInst, evalCode int) {

	log.Debugf("Link Eval\n")

	context.SetState(STATE_LINK_TRUE)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// State
const (
	STATE_NOT_STARTED int = 0

	STATE_LINK_FALSE int = 1
	STATE_LINK_TRUE  int = 2

	STATE_ENTERED int = 10
	STATE_READY   int = 20
	STATE_WAITING int = 30
	STATE_DONE    int = 40
)
