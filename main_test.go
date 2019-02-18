package main

import (
	"fmt"
	"testing"
)

//-----
func ExampleLoad() {
	fmt.Println("Load('PATH/TO/config.json')")
}

func TestLoad(t *testing.T) {
	configPath := "./config.json"
	if _, err := Load(configPath); err != nil {
		fmt.Println(" File config.json loaded")
	}
}

//-----
func ExampleParsePath() {
	fmt.Println("ParsePath('PATH/TO/QOS_HOME') to specific the '.qoscli' directory")
}

func TestParsePath(t *testing.T) {
	qosPath := "~/.qoscli"
	if _, err := ParsePath(qosPath); err != nil {
		fmt.Println("Find qoscli home directory")
	}
}

