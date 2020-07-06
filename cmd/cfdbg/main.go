package main

import (
	"flag"
	"github.com/PLAB-IO/aws-cloudformation-debugger/internal/cloudformation"
	"github.com/PLAB-IO/aws-cloudformation-debugger/internal/ui"
	awsCF "github.com/aws/aws-sdk-go/service/cloudformation"
	"log"
	"regexp"
	"strings"
	"time"
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
	events := lookupOriginalFailed(stackName)

	/*rows := [][]string {
		{"Timestamp", "Stack Name", "Fail Reason"},
	}
	for _, event := range events {
		rows = append(rows, []string{
			event.Timestamp.Format(time.Stamp),
			*event.StackName,
			*event.ResourceStatusReason,
		})
	}
	ui.PaintTable(rows)*/

	for _, event := range events {
		rows := [][]string {
			{"Timestamp", event.Timestamp.Format(time.RFC1123)},
			{"Stack Name", *event.StackName},
			{"Stack Status", *event.ResourceStatus},
			{"Logical Resource Id", *event.LogicalResourceId},
			{"FAILED Reason", *event.ResourceStatusReason},
		}
		ui.PaintTable(rows)
	}
}

func lookupOriginalFailed(stackName string) []awsCF.StackEvent {
	response :=  make([] awsCF.StackEvent, 0)
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

		if len(matched) == 2 {
			response = append(response, lookupOriginalFailed(matched[1])...)
			continue
		}

		response = append(response, event)
	}

	return response
}