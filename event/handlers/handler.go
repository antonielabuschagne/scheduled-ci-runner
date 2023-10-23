package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antonielabuschagne/event-triggered-jenkins-build/interfaces"
	"github.com/aws/aws-lambda-go/events"

	"github.com/joerdav/zapray"
)

type EventHandler struct {
	Log       *zapray.Logger
	JobRunner interfaces.IJobRunner
}

func NewEventHandler(log *zapray.Logger, j interfaces.IJobRunner) *EventHandler {
	return &EventHandler{
		Log:       log,
		JobRunner: j,
	}
}

func (eh EventHandler) Handle(ctx context.Context, e interface{}) (events.APIGatewayV2HTTPResponse, error) {
	eh.Log.Info("Handler triggered")
	eh.Log.Info(fmt.Sprintf("%+v", e))
	err := eh.JobRunner.RunBuild()
	if err != nil {
		eh.Log.Error(err.Error())
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	eh.Log.Info("Request completed")
	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusCreated}, nil
}