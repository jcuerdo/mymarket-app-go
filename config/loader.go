package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type Loader struct{
	Filename string
}

func (loader Loader) Load() Config {
	if loader.Filename != "" {
		return loadParams(loader.Filename)
	}

	return Config{}
}

func loadParams(fileName string) Config {
	config := Config{}
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("Error trying to read file %s", fileName))
	}

	if err = yaml.Unmarshal(contents, &config); err != nil {
		panic(fmt.Errorf("Config file %s has bad format", fileName))
	}

	return config
}