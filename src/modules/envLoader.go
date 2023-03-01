package modules

import (
	"crawling.com/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func LoadEnvironment() (*models.Environment, error) {

	f := ".secret/environment.yaml"
	buf, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("LoadEnvironment Get Error %v \n", err)
	}
	env := &models.Environment{}
	err = yaml.Unmarshal(buf, env)
	if err != nil {
		log.Fatalf("LoadEnvironment Unmarshal %v \n", err)
	}

	return env, nil

}
