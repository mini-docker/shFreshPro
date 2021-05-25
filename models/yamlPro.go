package models

import (
	"io/ioutil"
	"log"

	// "reflect"
	"gopkg.in/yaml.v2"
)

var Conf Yaml

func GetYaml() {
	// resultMap := make(map[string]interface{})
	conf := new(Yaml)
	yamlFile, err := ioutil.ReadFile("./conf/test.yaml")

	// conf := new(module.Yaml1)
	// yamlFile, err := ioutil.ReadFile("test.yaml")

	// conf := new(module.Yaml2)
	//  yamlFile, err := ioutil.ReadFile("test1.yaml")

	// log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	// log.Println("conf", conf, reflect.TypeOf(conf))
	Conf = *conf
	// log.Println("Confv", Conf)
}
