package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/colorstring"
)

func Check() (int, error) {
	var sopsFile string

	pushCommand := flag.NewFlagSet("check", flag.ExitOnError)
	data := AddStandardFlags(pushCommand, &sopsFile)
	if err := pushCommand.Parse(os.Args[2:]); err != nil {
		return 1, err
	}

	if &data == nil {
		log.Fatalf("Sops data missing. Did you specify -sops-file-path?")
	}

	ssmPathPrefix := data.Ssm["path-prefix"]
	if !strings.HasSuffix(ssmPathPrefix, "/") {
		log.Printf("Specified SSM prefix is missing the trailing /, adding")
		ssmPathPrefix = ssmPathPrefix + "/"
	}

	log.Printf("Checking SSM Path %s against local file %s\n", data.Ssm["path-prefix"], sopsFile)

	localParameters := GetSopsParameterData(data)
	existingParameters, err := GetExistingAWSParameters(client, ssmPathPrefix)
	if err != nil {
		return 10, err
	}

	localOnlyParameters, changedParameters, remoteOnlyParameters := BuildDiff(localParameters, existingParameters)

	writeDiff(localOnlyParameters, changedParameters, remoteOnlyParameters)
	return 0, nil
}

func writeDiff(localOnlyParameters []string, changedParameters []string, remoteOnlyParameters []string) {
	var sb strings.Builder
	var color = &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: false,
		Reset:   false,
	}
	if len(localOnlyParameters) == 0 {
		sb.WriteString(color.Color("[dark_gray]No New Parameters Found[reset]\n"))
	} else {
		for i := 0; i < len(localOnlyParameters); i++ {
			sb.WriteString(color.Color(fmt.Sprintf("[green]+[reset] new parameter %s\n", localOnlyParameters[i])))
		}
	}
	if len(changedParameters) == 0 {
		sb.WriteString(color.Color("[dark_gray]No Changed Parameters Found[reset]\n"))
	} else {
		for i := 0; i < len(changedParameters); i++ {
			sb.WriteString(color.Color(fmt.Sprintf("[yellow]~[reset] changed parameter %s\n", changedParameters[i])))
		}
	}

	if len(remoteOnlyParameters) == 0 {
		sb.WriteString(color.Color("[dark_gray]No Removed Parameters Found[reset]\n"))
	} else {
		for i := 0; i < len(remoteOnlyParameters); i++ {
			sb.WriteString(color.Color(fmt.Sprintf("[red]-[reset] removed parameter %s\n", remoteOnlyParameters[i])))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(color.Color(fmt.Sprintf("[green]%d[reset] parameters to add, [yellow]%d[reset] parameters to change, [red]%d[reset] parameters to remove", len(localOnlyParameters), len(changedParameters), len(remoteOnlyParameters))))

	fmt.Println(sb.String())
}
