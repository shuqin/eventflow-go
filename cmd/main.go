package main

import (
	"fmt"
	"github.com/qt/eventflow/pkg/compos"
	"github.com/qt/eventflow/pkg/define"
	"github.com/qt/eventflow/pkg/yaml"
	"reflect"
)


func main() {
	//testBasic()
	testComponentExecution()
}

func testBasic() {
	testPublicMemberConstruct()
	testPrivateMemberConstruct()
	testFactoryMethod()
	testEmptyObject()
	testSingleton()
	testClassMethod()
	testReadYaml()
}

func testPublicMemberConstruct() {
	result := &define.FlowResult{ 200 }
	fmt.Println(result)
	fmt.Println(reflect.TypeOf(result))
}

func testPrivateMemberConstruct() {
	//agent := &compos.AgentDetection{ "eventId", "agentId", 1 }
	//error: implicit assignment of unexported field 'eventId' in compos.AgentDetection literal

	empty := compos.AgentDetection{}
	fmt.Println(empty)
	fmt.Println(reflect.TypeOf(empty))

	detection := compos.AgentDetection{}.Instanceof("agentId", 1)
	fmt.Printf("d: %v\n", detection)

	detection2 := compos.AgentDetection{}.InstanceByNaming("agentId", 1)
	fmt.Printf("d2: %v\n", detection2)
	fmt.Println("eventId: " + detection2.GetEventId())

	detection3 := compos.NewAgentDetection("agentId", 1)
	fmt.Printf("d3: %v\n", detection3)
}

func testFactoryMethod() {
	adf := compos.GetAgentFactoryInstance()
	detection := adf.Instanceof("agentId", 2)
	fmt.Printf( "detection from factory method: %v\n", detection)
}

func testEmptyObject() {
	var emp compos.AgentDetection
	p1 := &emp
	demp := emp.Instanceof("agentId", 1)
	fmt.Printf("demp: %v\n", demp)
	fmt.Println(emp)
	fmt.Println(reflect.TypeOf(emp))

	var emp2 compos.AgentDetection
	p2 := &emp2

	fmt.Println("compare reference to empty agentDetection")
	compareObj(emp, emp2)


	fmt.Println("compare pointers to empty agentDetection")
	compareAgentDetection(p1, p2)

	emp3 := &compos.AgentDetection{}
	emp4 := &compos.AgentDetection{}
	fmt.Println("compare pointers to empty2 agentDetection")
	compareAgentDetection(emp3, emp4)
}

func compareObj(o1 interface{}, o2 interface{}) {
	fmt.Printf("o1: %v\n", o1)
	fmt.Printf("o2: %v\n", o2)
	fmt.Printf("o1-p: %p\n", o1)
	fmt.Printf("o2-p: %p\n", o2)
	fmt.Printf("&o1: %v\n", &o1)
	fmt.Printf("&o2: %v\n", &o2)
	fmt.Printf("o1 == o2: %v\n", (o1 == o2))
	fmt.Printf("&o1 == &o2: %v\n", &o1 == &o2)
}

func compareAgentDetection(adp1 *compos.AgentDetection, adp2 *compos.AgentDetection) {
	fmt.Printf("adp1: %v\n", adp1)
	fmt.Printf("adp2: %v\n", adp2)
	fmt.Printf("adp1-p: %p\n", adp1)
	fmt.Printf("adp2-p: %p\n", adp2)
	fmt.Printf("*adp1: %v\n", *adp1)
	fmt.Printf("*adp2: %v\n", *adp2)
	fmt.Printf("&adp1: %v\n", &adp1)
	fmt.Printf("&adp2: %v\n", &adp2)
	fmt.Printf("adp1 == adp2: %v\n", (adp1 == adp2))
	fmt.Printf("&adp1 == &adp2: %v\n", &adp1 == &adp2)
	fmt.Printf("*adp1 == *adp2: %v\n", *adp1 == *adp2)
}

func testSingleton() {
	adf := compos.GetAgentFactoryInstance()
	adf2 := compos.GetAgentFactoryInstance()
	compare(adf, adf2)
}

func compare(p1 *compos.AgentDetectionFactory, p2 *compos.AgentDetectionFactory) {
	fmt.Printf("adf: %v\n", p1)
	fmt.Printf("adf2: %v\n", p2)
	fmt.Printf("adf-p: %p\n", p1)
	fmt.Printf("adf2-p: %p\n", p2)
	fmt.Printf("*adf: %v\n", *p1)
	fmt.Printf("*adf2: %v\n", *p2)
	fmt.Printf("&adf: %v\n", &p1)
	fmt.Printf("&adf2: %v\n", &p2)
	fmt.Printf("adf == adf2: %v\n", (p1 == p2))
	fmt.Printf("&adf == &adf2: %v\n", &p1 == &p2)
	fmt.Printf("*adf == *adf2: %v\n", *p1 == *p2)
}

func testClassMethod() {
	withoutStar := compos.AgentDetection{}
	ad := withoutStar.Instanceof("agentId111", 1)
	fmt.Printf("ad: %v\n", ad)

	withoutStar2 := &compos.AgentDetection{}
	ad2 := withoutStar2.Instanceof("agentId222", 1)
	fmt.Printf("ad: %v\n", ad2)

	withStar := compos.AgentDetectionFactory{}
	ad3 := withStar.Instanceof("agentId333", 1)
	fmt.Printf("ad3: %v\n", ad3)

	withStar2 := &compos.AgentDetectionFactory{}
	ad4 := withStar2.Instanceof("agentId444", 1)
	fmt.Printf("ad4: %v\n", ad4)
}

func testReadYaml() {
	yaml.ReadYaml()
}

func testComponentExecution() {

	//ctx := &compos.Context{ &event } //Cannot assign a value to the unexported field 'event'

	detection := compos.AgentDetection{}.Instanceof("agentId", 1)
	var event define.EventData = detection
	ctx := compos.Context{}.New(&event)
	cf := compos.GetComponentFactoryInstance()
	cf.Init()

	ch := make(chan int, 4)
	defer close(ch)

	go func() {
		ch <- 0
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- -1
	}()

	for {
		select {
		case msg := <- ch:
			if (msg == 0) {
				execDirect(ctx)
			}
			if (msg == 1) {
				execRange(ctx)
			}
			if (msg == 2) {
				execRange2(ctx, cf)
			}
			if (msg == 3) {
				execByArrange(ctx, cf)
			}
			if (msg == -1) {
				return
			}
		}
	}
}


func execDirect(ctx *compos.Context) {
	fmt.Println("execDirect")
	fc := &compos.HostCompletionComponent{}
	saver := &compos.DetectionSaver{}
	fc.Process(ctx)
	saver.Process(ctx)
}

func execRange(ctx *compos.Context)  {
	fmt.Println("execRange")
	components := [...]compos.FlowComponent{&compos.HostCompletionComponent{}, &compos.DetectionSaver{} }
	for _, comp := range components {
		res, err := comp.Process(ctx)
		if err != nil {
			fmt.Println(err)
		}
		if res.Code == define.TERMINAL {
			panic("Flow terminated")
		}
	}
}

func execRange2(ctx *compos.Context, cf *compos.ComponentFactory)  {
	fmt.Println("execRange2")
	components := compos.MapValue([]int { compos.HostCompletionComponentName, compos.DetectionSaverName}, cf.Instanceof)
	compos.Exec(components, ctx)
}

func execByArrange(ctx *compos.Context, cf *compos.ComponentFactory) {
	fmt.Println("execByArrange")
	flow := yaml.Read("./configs/eventflow.yml")
	components := make([]compos.FlowComponent, 2)
	for _, f := range flow.Flows {
		fmt.Println(f.BizTypes)
		for i, c := range f.ComponentConfigs {
			if trueComp := cf.GetByName(c.Name) ; trueComp != nil {
				components[i] = trueComp
			}
		}
	}

	compos.Exec(components, ctx)
}


