package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var (
	Region string
	awsSession *session.Session
)

func GetFailEvents(stackName string) []cloudformation.StackEvent {
	stackEventsInput := cloudformation.DescribeStackEventsInput{
		StackName: &stackName,
	}
	svc := cloudformation.New(awsSession, &aws.Config{
		Region: &Region,
	})
	events := make([]cloudformation.StackEvent,0)

	for {
		outputs, err := svc.DescribeStackEvents(&stackEventsInput)
		if err != nil {
			panic(err)
		}

		for _, event := range outputs.StackEvents {
			events = append(events, *event)
		}

		if outputs.NextToken == nil {
			break
		}
		stackEventsInput.SetNextToken(*outputs.NextToken)
	}

	return events
}

func SetProfile(profile string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
	})
	if err != nil {
		return err
	}
	awsSession = sess

	return nil
}
