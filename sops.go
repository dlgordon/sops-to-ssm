package main

import (
	"go.mozilla.org/sops/v3/decrypt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type sopsData struct {
	Environment map[string]string `yaml:"environment"`
	Secrets     map[string]string `yaml:"secrets"`
}

func LoadSopsFile() (sopsData, error) {
	b, err := os.ReadFile(sopsFile)
	if err != nil {
		log.Printf("Couldnt read file %v", sopsFile)
		return sopsData{}, err
	}

	decryptedFile, err := decrypt.Data(b, "yaml")
	if err != nil {
		log.Printf("Couldnt read decrypt sops file %v. Error %v", sopsFile, err)
		return sopsData{}, err
	}

	env := sopsData{}

	if err := yaml.Unmarshal(decryptedFile, &env); err != nil {
		log.Printf("Couldnt marhsal file %v", err)
		return sopsData{}, err
	}

	return env, nil
}
