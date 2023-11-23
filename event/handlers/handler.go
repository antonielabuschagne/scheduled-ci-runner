package handlers

import (
	"context"
	"net/http"

	"github.com/antonielabuschagne/event-triggered-jenkins-build/interfaces"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"

	"github.com/joerdav/zapray"
)

type EventHandler struct {
	Log        *zapray.Logger
	JobBuilder interfaces.JobBuilder
}

func NewEventHandler(log *zapray.Logger, j interfaces.JobBuilder) *EventHandler {
	return &EventHandler{
		Log:        log,
		JobBuilder: j,
	}
}

func (eh EventHandler) Handle(ctx context.Context, e interface{}) (events.APIGatewayV2HTTPResponse, error) {
	eh.Log.Info("Handler triggered", zap.Any("event", e))
	err := eh.JobBuilder.Build()
	if err != nil {
		eh.Log.Error(err.Error())
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	eh.Log.Info("Request completed")
	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusCreated}, nil
}
