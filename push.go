package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func Push() (int, error) {
	var sopsFile string

	pushCommand := flag.NewFlagSet("push", flag.ExitOnError)
	shouldDelete := pushCommand.Bool("remove-missing", false, "Remove parameters that are mising from the local file")
	data := AddStandardFlags(pushCommand, &sopsFile)
	if err := pushCommand.Parse(os.Args[2:]); err != nil {
		return 1, err
	}

	if data == nil {
		log.Fatalf("Sops data missing. Did you specify -sops-file-path?")
	}

	ssmPathPrefix := data.Ssm["path-prefix"]
	if !strings.HasSuffix(ssmPathPrefix, "/") {
		log.Printf("Specified SSM prefix is missing the trailing /, adding")
		ssmPathPrefix = ssmPathPrefix + "/"
	}

	localParameters := GetSopsParameterData(data)
	existingParameters, err := GetExistingAWSParameters(client, ssmPathPrefix)
	if err != nil {
		return 10, err
	}

	localOnlyParameters, changedParameters, remoteOnlyParameters := BuildDiff(localParameters, existingParameters)
	if len(localOnlyParameters) > 0 {
		var ssmCreates []*ssm.PutParameterInput
		for i := 0; i < len(localOnlyParameters); i++ {
			parameterName := localOnlyParameters[i]
			parameterValue := localParameters[parameterName].Value
			parameterType := localParameters[parameterName].Type
			ssmName := strings.Join([]string{ssmPathPrefix, parameterName}, "")
			ssmCreates = append(ssmCreates, &ssm.PutParameterInput{Name: &ssmName, Value: &parameterValue, Type: types.ParameterType(parameterType), Overwrite: false})
		}
		PutParameters(client, ssmCreates)
	}

	if len(changedParameters) > 0 {
		var ssmUpdates []*ssm.PutParameterInput
		for i := 0; i < len(changedParameters); i++ {
			parameterName := changedParameters[i]
			parameterValue := localParameters[parameterName].Value
			parameterType := localParameters[parameterName].Type
			ssmName := strings.Join([]string{ssmPathPrefix, parameterName}, "")
			ssmUpdates = append(ssmUpdates, &ssm.PutParameterInput{Name: &ssmName, Value: &parameterValue, Type: types.ParameterType(parameterType), Overwrite: true})
		}
		PutParameters(client, ssmUpdates)
	}

	if *shouldDelete && len(remoteOnlyParameters) > 0 {
		var ssmDeletes []string
		for i := 0; i < len(remoteOnlyParameters); i++ {
			ssmDeletes = append(ssmDeletes, strings.Join([]string{ssmPathPrefix, remoteOnlyParameters[i]}, ""))
		}
		DeleteParameters(client, &ssm.DeleteParametersInput{Names: ssmDeletes})
	}

	return 0, nil
}
