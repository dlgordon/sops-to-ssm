package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var client *ssm.Client

func main() {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error getting AWS configuration, %v", err)
		os.Exit(2)
	}
	client = ssm.NewFromConfig(awsConfig)
	switch os.Args[1] {
	case "push":
		Push()
	case "check":
		Check()
	default:
		log.Fatal("expected 'push' or 'check' subcommands")
	}
}
