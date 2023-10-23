package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/joho/godotenv"
)

type CDKStackProps struct {
	JenkinsApiUser       string
	JenkinsApiToken      string
	JenkinsBuildEndpoint string
	StackProps           awscdk.StackProps
}

func NewCDKStack(scope constructs.Construct, cdkProps CDKStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, aws.String("event-triggered-build"), &cdkProps.StackProps)

	vpc := awsec2.NewVpc(stack, aws.String("EventTriggerBuild"), &awsec2.VpcProps{
		MaxAzs: aws.Float64(2),
	})

	eventHandler := awslambdago.NewGoFunction(stack, aws.String("scheduledEventHandler"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        aws.String("../event/handlers/scheduledevent"),
		Bundling: &awslambdago.BundlingOptions{
			GoBuildFlags: &[]*string{aws.String(`-ldflags "-s -w" -tags lambda.norpc`)},
		},
		Environment: &map[string]*string{
			"JENKINS_API_USER":     &cdkProps.JenkinsApiUser,
			"JENKINS_API_TOKEN":    &cdkProps.JenkinsApiToken,
			"JENKINS_JOB_ENDPOINT": &cdkProps.JenkinsBuildEndpoint,
		},
		MemorySize: aws.Float64(1024),
		Tracing:    awslambda.Tracing_ACTIVE,
		Timeout:    awscdk.Duration_Millis(aws.Float64(300000)),
		Vpc:        vpc,
	})

	// EventBridget scheduled rule
	awsevents.NewRule(stack, aws.String("ScheduledTask"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Rate(awscdk.Duration_Millis(aws.Float64(300000))),
		Targets: &[]awsevents.IRuleTarget{
			awseventstargets.NewLambdaFunction(eventHandler, &awseventstargets.LambdaFunctionProps{}),
		},
	})
	return stack
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		println("AWS_REGION undefined")
	}
	awsAccount := os.Getenv("AWS_ACCOUNT_ID")
	if awsAccount == "" {
		panic("AWS_ACCOUNT_ID undefined")
	}
	jenkinsApiUser := os.Getenv("JENKINS_API_USER")
	if jenkinsApiUser == "" {
		panic("JENKINS_API_USER undefined")
	}
	jenkinsBuildEndpoint := os.Getenv("JENKINS_JOB_ENDPOINT")
	if jenkinsBuildEndpoint == "" {
		panic("JENKINS_JOB_ENDPOINT undefined")
	}
	jenkinsApiToken, err := GetJenkinsApiToken(awsRegion)
	if err != nil {
		panic(fmt.Sprintf("unable to get JENKINS_API_TOKEN secret: %v", err.Error()))
	}

	app := awscdk.NewApp(nil)
	NewCDKStack(app, CDKStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region:  aws.String(awsRegion),
				Account: aws.String(awsAccount),
			},
		},
		JenkinsApiUser:       jenkinsApiUser,
		JenkinsBuildEndpoint: jenkinsBuildEndpoint,
		JenkinsApiToken:      jenkinsApiToken,
	})
	app.Synth(nil)
}

func GetJenkinsApiToken(region string) (secret string, err error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return
	}

	svc := secretsmanager.NewFromConfig(config)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("JENKINS_API_TOKEN"),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return
	}
	// Decrypts secret using the associated KMS key.
	secretData := *result.SecretString
	// Unmarshal into data structure.
	var sec JenkinsApiTokenData
	err = json.Unmarshal([]byte(secretData), &sec)
	if err != nil {
		return
	}
	secret = sec.JenkinsApiToken
	return
}

type JenkinsApiTokenData struct {
	JenkinsApiToken string `json:"JENKINS_API_TOKEN"`
}
