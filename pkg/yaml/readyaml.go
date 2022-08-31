package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type CommonConfig struct {
	eventType string `yaml:"eventType"`
	EventSourceType string `yaml:"eventSourceType"`
	Model string `yaml:"model"`
	Way string `yaml:"way"`
	OriginParamType string `yaml:"originParamType"`
	BuilderType string `yaml:"builderType"`
	ComponentParamType string `yaml:"componentParamType"`
}

type ComponentConfig struct {
	Name string   `yaml:"name"`
}

type Flow struct {
	BizTypes         []string       `yaml:"bizTypes"`
	CommonConfig     CommonConfig   `yaml:"commonConfig"`
	ComponentConfigs []ComponentConfig `yaml:"componentConfigs"`
}

type EventFlow struct {
	Version string `yaml:"version"`
	Flows []Flow     `yaml:"flows"`
}

func Read(filename string) EventFlow {
	ef := EventFlow{}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(string(yamlFile))

	err2 := yaml.Unmarshal(yamlFile, &ef)
	if err2 != nil {
		log.Fatalf("error: %v", err2)
	}
	return ef
}

func ReadYaml() {

	ef := Read("./configs/eventflow.yml")
	fmt.Println(ef)
	for _, f := range ef.Flows {
		fmt.Println(f.BizTypes)
		fmt.Println(f.CommonConfig)
		for _, c := range f.ComponentConfigs {
			fmt.Println(c.Name)
		}
	}
}
