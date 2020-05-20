package main

import "flag"

var yamlFileFlag = flag.String("yaml", "", "YAML file configuration to be loaded")
var jsonFileFlag = flag.String("json", "", "JSON file configuration to be loaded")

func init() {
	flag.Parse()
}
