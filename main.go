package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var (
	ssmPathPrefix string
	sopsFile      string
)

func init() {
	flag.StringVar(&ssmPathPrefix, "path-prefix", "", "Specify the path to search SSM from")
	flag.StringVar(&sopsFile, "sops-file-path", "", "Specific a path to read your local env files")
	flag.Parse()

	if !strings.HasSuffix(ssmPathPrefix, "/") {
		log.Printf("Specified SSM prefix is missing the trailing /, adding")
		ssmPathPrefix = ssmPathPrefix + "/"
	}
}

func main() {
	if "" == sopsFile {
		log.Fatal("sops file not specified")
		os.Exit(1)
	}

	if "/" == ssmPathPrefix {
		log.Print("SSM path not specified or root path specified")
		os.Exit(1)
	}
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error getting AWS configuration, %v", err)
		os.Exit(2)
	}

	sopsFileParameters, err := LoadSopsFile()
	if err != nil {
		log.Printf("Couldnt decrypt sops file, %v", err)
		os.Exit(3)
	}

	client := ssm.NewFromConfig(awsConfig)
	existingSsmParameters := GetParameters(client)

	newParameters, changedParameters := BuildDiff(sopsFileParameters.Environment, existingSsmParameters, false)
	log.Printf("Found %d New Environment Parameters", len(newParameters))
	log.Printf("Found %d Changed Environment Parameters", len(changedParameters))

	newSecrets, changedSecrets := BuildDiff(sopsFileParameters.Secrets, existingSsmParameters, true)
	log.Printf("Found %d New Secret Parameters", len(newSecrets))
	log.Printf("Found %d Changed Secret Parameteres", len(changedSecrets))

	if len(newParameters) > 0 {
		log.Print("Creating new parameters")
		PutParameters(client, newParameters)
	}

	if len(changedParameters) > 0 {
		log.Print("Updating changed parameters")
		PutParameters(client, changedParameters)
	}

	if len(newSecrets) > 0 {
		log.Print("Creating new secrets")
		PutParameters(client, newSecrets)
	}

	if len(changedSecrets) > 0 {
		log.Print("Updating changed secrets")
		PutParameters(client, changedSecrets)
	}

}
