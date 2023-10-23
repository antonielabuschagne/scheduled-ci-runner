package main

import (
	"os"

	"github.com/antonielabuschagne/event-triggered-jenkins-build/event/handlers"
	"github.com/antonielabuschagne/event-triggered-jenkins-build/jenkins"
	"github.com/aws/aws-lambda-go/lambda"

	"go.uber.org/zap"

	"github.com/joerdav/zapray"
)

func main() {
	log, err := zapray.NewProduction()
	if err != nil {
		log.Fatal("failed to create logger", zap.String("error", err.Error()))
	}
	jenkinsApiUser := os.Getenv("JENKINS_API_USER")
	if jenkinsApiUser == "" {
		log.Fatal("missing JENKINS_API_USER")
	}
	jenkinsApiToken := os.Getenv("JENKINS_API_TOKEN")
	if jenkinsApiToken == "" {
		log.Fatal("missing JENKINS_API_TOKEN")
	}
	jenkinsBuildEndpoint := os.Getenv("JENKINS_JOB_ENDPOINT")
	if jenkinsBuildEndpoint == "" {
		log.Fatal("missing JENKINS_JOB_ENDPOINT")
	}

	jjb := jenkins.NewJobRunner(log, jenkinsApiUser, jenkinsApiToken, jenkinsBuildEndpoint)
	eh := handlers.NewEventHandler(log, jjb)

	lambda.Start(eh.Handle)
}
