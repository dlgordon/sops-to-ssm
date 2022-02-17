package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetExistingAWSParameters(client *ssm.Client, ssmPathPrefix string) (map[string]parameterData, error) {
	existingAWSParameters := make(map[string]parameterData)

	var nextToken *string
	for {
		params, err := client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{Path: &ssmPathPrefix, Recursive: true, WithDecryption: true, NextToken: nextToken})
		if err != nil {
			log.Fatalf("failed to get parameters, %v", err)
		}

		for i := 0; i < len(params.Parameters); i++ {
			// Prune the prefix and store
			p := params.Parameters[i]
			existingAWSParameters[strings.TrimPrefix(*p.Name, ssmPathPrefix)] = parameterData{Name: *p.Name, Value: *p.Value, Type: string(p.Type)}
			// log.Printf("Parameter: %s; Value: %s", *params.Parameters[i].Name, *params.Parameters[i].Value)
		}

		if params.NextToken == nil {
			break
		}
		//log.Printf("Log %d parameters. next token: %s", len(params.Parameters), *params.NextToken)
		nextToken = params.NextToken
	}
	return existingAWSParameters, nil
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

func DeleteParameters(client *ssm.Client, parameters *ssm.DeleteParametersInput) {
	_, err := client.DeleteParameters(context.TODO(), parameters)
	if err != nil {
		log.Fatalf("Error setting a parameter: %v", err)
	}
}
