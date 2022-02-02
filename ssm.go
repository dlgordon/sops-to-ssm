package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type awsParameter struct {
	Path  string
	Value string
	Type  string
}

func GetParameters(client *ssm.Client) map[string]awsParameter {
	existingAWSParameters := make(map[string]awsParameter)

	var nextToken *string
	for {
		params, err := client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{Path: &ssmPathPrefix, Recursive: true, WithDecryption: true, NextToken: nextToken})
		if err != nil {
			log.Fatalf("failed to get parameters, %v", err)
		}

		for i := 0; i < len(params.Parameters); i++ {
			// Prune the prefix and store
			p := params.Parameters[i]
			existingAWSParameters[strings.TrimPrefix(*p.Name, ssmPathPrefix)] = awsParameter{Path: *p.Name, Value: *p.Value, Type: string(p.Type)}
			// log.Printf("Parameter: %s; Value: %s", *params.Parameters[i].Name, *params.Parameters[i].Value)
		}

		if params.NextToken == nil {
			break
		}
		//log.Printf("Log %d parameters. next token: %s", len(params.Parameters), *params.NextToken)
		nextToken = params.NextToken
	}
	return existingAWSParameters
}

func BuildDiff(localVariables map[string]string, existingVariables map[string]awsParameter, secureString bool) ([]*ssm.PutParameterInput, []*ssm.PutParameterInput) {
	var newParameters []*ssm.PutParameterInput
	var changedParameters []*ssm.PutParameterInput

	for k, v := range localVariables {
		currentValue, ok := existingVariables[k]
		if !ok {
			var newPath = strings.Join([]string{ssmPathPrefix, k}, "")
			var newValue = v
			var parameterType = types.ParameterTypeString
			if secureString {
				parameterType = types.ParameterTypeSecureString
			}
			newParameters = append(newParameters, &ssm.PutParameterInput{
				Name:  &newPath,
				Value: &newValue,
				Type:  parameterType,
			})
			continue
		}

		//log.Printf("Local Key: %s Value: %s vs Remote value: %s", k, v, currentValue.Value)
		if strings.Compare(v, currentValue.Value) != 0 {
			var newPath = strings.Join([]string{ssmPathPrefix, k}, "")
			var newValue = v
			changedParameters = append(changedParameters, &ssm.PutParameterInput{
				Name:      &newPath,
				Value:     &newValue,
				Type:      types.ParameterType(currentValue.Type),
				Overwrite: true,
			})
			//log.Printf("Name %s Value %s", *changedParameters[len(changedParameters)-1].Name, *changedParameters[len(changedParameters)-1].Value)
		}

	}

	return newParameters, changedParameters
}

func PutParameters(client *ssm.Client, parameters []*ssm.PutParameterInput) {
	for i := 0; i < len(parameters); i++ {
		_, err := client.PutParameter(context.TODO(), parameters[i])
		if err != nil {
			log.Fatalf("Error setting a parameter: %v", err)
		}
		//log.Printf("Name: %s Value: %s", *parameters[i].Name, *parameters[i].Value)
	}
}
