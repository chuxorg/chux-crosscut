package secrets

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func SetEnvironment(functionName, region string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Error creating AWS session: %v", err)
	}

	lambdaClient := lambda.New(sess)

	input := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
		Environment: &lambda.Environment{
			Variables: map[string]*string{
				"AWS_REGION": aws.String(region),
			},
		},
	}

	_, err = lambdaClient.UpdateFunctionConfiguration(input)
	if err != nil {
		log.Fatalf("Error updating Lambda function configuration: %v", err)
	}

	log.Printf("Environment variable set for Lambda function: %s", functionName)
}
