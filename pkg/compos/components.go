package compos

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/qt/eventflow/pkg/define"
)

type AgentDetection struct {
	eventId string
	agentId string
	osType  int
}

func getRandomInt() int {
	nanotime := int64(time.Now().Nanosecond())
	rand.Seed(nanotime)
	return rand.Int()
}

func (AgentDetection) Instanceof(agentId string, osType int) *AgentDetection {
	eventId := strconv.Itoa(int((time.Now().Unix()))) + strconv.Itoa(getRandomInt())
	return &AgentDetection{ eventId, agentId, osType}
}

func (AgentDetection) InstanceByNaming(agentId string, osType int) *AgentDetection {
	return &AgentDetection{ agentId:agentId, osType: osType}
}

func NewAgentDetection(agentId string, osType int) *AgentDetection {
	eventId := strconv.Itoa(int((time.Now().Unix()))) + strconv.Itoa(getRandomInt())
	return &AgentDetection{ eventId, agentId, osType}
}

func (detection *AgentDetection) GetData() interface{} {
	return detection
}

func (detection *AgentDetection) GetType() string {
	return "Detection"
}

func (detection *AgentDetection) GetEventId() string {
	return detection.eventId
}

func (detection *AgentDetection) GetAgentId() string {
	return detection.agentId
}

func (detection *AgentDetection) SetOsType(osType int) (err error) {
	if osType != 1 && osType != 2 {
		return &define.IllegalArgumentError{"invalid parameter", nil}
	}
	(*detection).osType = osType
	return nil
}


type AgentDetectionFactory struct {
}

func (*AgentDetectionFactory) Instanceof(agentId string, osType int) *AgentDetection {
	eventId := strconv.Itoa(int((time.Now().Unix()))) + strconv.Itoa(getRandomInt())
	return &AgentDetection{ eventId, agentId, osType}
}

func (*AgentDetectionFactory) InstanceByNaming(agentId string, osType int) *AgentDetection {
	return &AgentDetection{ agentId:agentId, osType: osType}
}

var (
    facLock     *sync.Mutex = &sync.Mutex{}
    adf *AgentDetectionFactory
)

func GetAgentFactoryInstance() *AgentDetectionFactory {
	if adf == nil {
		facLock.Lock()
		defer facLock.Unlock()
		if adf == nil {
			adf = &AgentDetectionFactory{}
			fmt.Println("AgentDetectionFactory single instance...")
		}
	}
	return adf
}

type Context struct {
	event *define.EventData
}

func (ctx Context) New(e *define.EventData) *Context  {
	return &Context{ e }
}

// Do not need declared explicitly
//func NewContext(event *EventData) *Context {
//	return &Context{event }
//}

func (ctx *Context) GetEventData() *define.EventData {
	return ctx.event
}

const (
	HostCompletionComponentName = iota
	DetectionSaverName
)

type FlowComponent interface {
	GetName() string
	Process(context *Context) (result define.FlowResult, error error)
}

func getName(o interface{}) string  {
	return fmt.Sprintf("%v", reflect.TypeOf(o))
}

type HostCompletionComponent struct {

}

func (hcc HostCompletionComponent) GetName() string  {
	return getName(hcc)
}

func (hcc HostCompletionComponent) Process(context *Context) (result define.FlowResult, error error) {
	event := context.GetEventData()
	err := (*event).SetOsType(1)
	if err != nil {
		return define.FlowResult{define.TERMINAL}, err
	}
	fmt.Println(*event)
	return define.FlowResult{define.CONTINUE}, nil
}

type DetectionSaver struct {
}

func (saver DetectionSaver) Process(ctx *Context) (result define.FlowResult, error error) {
	event := ctx.GetEventData()
	fmt.Printf("save detection. eventId: %s\n", (*event).GetEventId())
	return define.FlowResult{200}, nil
}

func (saver DetectionSaver) GetName() string  {
	return getName(saver)
}

type Empty struct {

}

func (Empty) Process(ctx *Context) (result define.FlowResult, error error) {
	return define.FlowResult{200}, nil
}

func (empty Empty) GetName() string  {
	return getName(empty)
}

var (
	lock     *sync.Mutex = &sync.Mutex{}
	hostComponent *HostCompletionComponent
	detectionSaverComponent *DetectionSaver
	empty *Empty
	componentFactory *ComponentFactory
)

func GetHostCompletionComponentInstance() *HostCompletionComponent {
	if hostComponent == nil {
		lock.Lock()
		defer lock.Unlock()
		if hostComponent == nil {
			hostComponent = &HostCompletionComponent{}
			fmt.Println("HostCompletionComponent single instance...")
		}
	}
	return hostComponent
}

func GetDetectionSaverComponentInstance() *DetectionSaver {
	if detectionSaverComponent == nil {
		lock.Lock()
		defer lock.Unlock()
		if detectionSaverComponent == nil {
			detectionSaverComponent = &DetectionSaver{}
			fmt.Println("DetectionSaver single instance...")
		}
	}
	return detectionSaverComponent
}

func GetEmptyInstance() *Empty {
	if empty == nil {
		lock.Lock()
		defer lock.Unlock()
		if empty == nil {
			empty = &Empty{}
			fmt.Println("Empty single instance...")
		}
	}
	return empty
}

type ComponentFactory struct {

}

func GetComponentFactoryInstance() *ComponentFactory {
	if componentFactory == nil {
		lock.Lock()
		defer lock.Unlock()
		if componentFactory == nil {
			componentFactory = &ComponentFactory{}
			fmt.Println("ComponentFactory single instance...")
		}
	}
	return componentFactory
}

func (*ComponentFactory) Instanceof(compName int) FlowComponent {
	switch compName {
	case HostCompletionComponentName:
		return GetHostCompletionComponentInstance()
	case DetectionSaverName:
		return GetDetectionSaverComponentInstance()
	default:
		return GetEmptyInstance()
	}
}

var componentMap = map[string]FlowComponent {}

func (*ComponentFactory) Init() {
	hcc := GetHostCompletionComponentInstance()
	componentMap[hcc.GetName()] = hcc
	saver := GetDetectionSaverComponentInstance()
	componentMap[saver.GetName()] = saver
}

func (*ComponentFactory) GetByName(name string) FlowComponent  {
	return componentMap[name]
}


func MapValue(ints []int, f func(i int) FlowComponent) []FlowComponent  {
	result := make([]FlowComponent, len(ints))
	for i, v := range ints {
		result[i] = f(v)
	}
	return result
}

func Exec(components []FlowComponent, ctx *Context) (err error) {
	for _, comp := range components {
		res, err := comp.Process(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if res.Code == define.TERMINAL {
			panic("Flow terminated")
		}
	}
	return nil
}
