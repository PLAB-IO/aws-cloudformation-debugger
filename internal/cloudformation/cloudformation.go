package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"regexp"
)

var (
	Region string
	awsSession *session.Session
	completeStatus = []string{
		"CREATE_COMPLETE",
		"UPDATE_COMPLETE",
		"UPDATE_ROLLBACK_COMPLETE",
	}
)

func GetFailEvents(stackName string) []cloudformation.StackEvent {
	stackEventsInput := cloudformation.DescribeStackEventsInput{
		StackName: &stackName,
	}
	svc := cloudformation.New(awsSession, &aws.Config{
		Region: &Region,
	})
	events := make([]cloudformation.StackEvent,0)

	// BECAUSE NESTED STACK ARE ARN AND NOT STACK NAME / WTF AWS ???
	re := regexp.MustCompile(`arn:aws:cloudformation:.*/(.*)/.*`)
	matched := re.FindStringSubmatch(stackName)
	realStackName := stackName
	if len(matched) == 2 {
		realStackName = matched[1]
	}

	for {
		outputs, err := svc.DescribeStackEvents(&stackEventsInput)
		if err != nil {
			panic(err)
		}
		breakLoop := false

		for _, event := range outputs.StackEvents {

			// MUST STOP AT NEXT COMPLETE EVENT THAT REFER THE LOGICAL STACK NAME
			if len(events) > 0 &&
							*event.LogicalResourceId == realStackName &&
							contains(completeStatus, *event.ResourceStatus) {
				breakLoop = true
				break
			}
			events = append(events, *event)
		}

		if outputs.NextToken == nil || breakLoop {
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

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}