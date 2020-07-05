package main

import (
	"flag"
	"github.com/PLAB-IO/aws-cloudformation-debugger/internal/cloudformation"
	"log"
	"regexp"
	"strings"
)

var (
	stackName = ""
)

func main() {
	profile := flag.String("profile", "default", "AWS Profile")
	region := flag.String("region", "eu-west-1", "AWS Region")
	_stackName := flag.String("stack-name", "", "Cloudformation Stack name")
	flag.Parse()

	if *profile == "" {
		log.Fatal("Please provide --profile option")
	}

	if *region == "" {
		log.Fatal("Please provide --region option")
	}

	if *_stackName == "" {
		log.Fatal("Please provide --stack-name option")
	}

	stackName = *_stackName
	cloudformation.Region = *region

	if err := cloudformation.SetProfile(*profile); err != nil {
		panic(err)
	}
	lookupOriginalFailed(stackName)
}

func lookupOriginalFailed(stackName string) {
	events := cloudformation.GetFailEvents(stackName)
	for _, event := range events {
		if !strings.Contains(*event.ResourceStatus, "FAILED") {
			continue
		}
		if "Resource creation cancelled" == *event.ResourceStatusReason {
			continue
		}

		re := regexp.MustCompile(`Embedded stack (arn:.*) was not successfully`)
		matched := re.FindStringSubmatch(*event.ResourceStatusReason)
		//println(matched[1])

		//println("=====================")

		if len(matched) == 2 {
			lookupOriginalFailed(matched[1])
			continue
		}

		println(event.String())
	}
}