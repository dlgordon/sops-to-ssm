package main

import (
	"log"
	"os"
	"strings"

	"go.mozilla.org/sops/v3/decrypt"
	"gopkg.in/yaml.v3"
)

type sopsData struct {
	Ssm         map[string]string `yaml:"ssm"`
	Environment map[string]string `yaml:"environment"`
	Secrets     map[string]string `yaml:"secrets"`
}

type parameterData struct {
	Name  string
	Value string
	Type  string
}

func LoadSopsFile(sopsFile string) (sopsData, error) {
	b, err := os.ReadFile(sopsFile)
	if err != nil {
		log.Printf("Couldnt read file %v", sopsFile)
		return sopsData{}, err
	}

	decryptedFile, err := decrypt.Data(b, "yaml")
	if err != nil {
		log.Printf("Couldn't decrypt sops file %v. Error %v", sopsFile, err)
		return sopsData{}, err
	}

	env := sopsData{}

	if err := yaml.Unmarshal(decryptedFile, &env); err != nil {
		log.Printf("File structure is not valid: %v", err)
		return sopsData{}, err
	}

	return env, nil
}

func GetSopsParameterData(sopsData *sopsData) map[string]parameterData {
	sopsParameters := make(map[string]parameterData)
	for k, v := range sopsData.Environment {
		sopsParameters[k] = parameterData{Name: k, Value: v, Type: "String"}
	}
	for k, v := range sopsData.Secrets {
		sopsParameters[k] = parameterData{Name: k, Value: v, Type: "SecureString"}
	}
	return sopsParameters
}

func BuildDiff(localVariables map[string]parameterData, remoteVariables map[string]parameterData) (
	[]string, []string, []string) {
	var changedParameters []string

	localOnlyParameters := keysMissing(localVariables, remoteVariables)
	remoteOnlyParameters := keysMissing(remoteVariables, localVariables)

	for k, v := range localVariables {
		currentValue, ok := remoteVariables[k]
		if !ok {
			continue
		}

		if strings.Compare(v.Value, currentValue.Value) == 0 && strings.Compare(v.Type, currentValue.Type) == 0 {
			continue
		} else {
			changedParameters = append(changedParameters, k)
		}
	}

	return localOnlyParameters, changedParameters, remoteOnlyParameters
}

func keysMissing(source map[string]parameterData, destination map[string]parameterData) []string {
	var keysMissing []string
	for k, _ := range source {
		_, ok := destination[k]
		if !ok {
			keysMissing = append(keysMissing, k)
		}
	}
	return keysMissing
}
