package jenkins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joerdav/zapray"
)

func TestRunBuild(t *testing.T) {
	logger, err := zapray.NewDevelopment()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		statusCode    int
		expectedError string
		user          string
		token         string
		description   string
	}{
		{
			description: "given a 201 response, build started successfully",
			statusCode:  http.StatusCreated,
		},
		{
			description:   "given, a 40X response, build failed",
			statusCode:    http.StatusBadRequest,
			expectedError: "unable to build job - status: 400 body: ",
		},
		{
			description:   "given, a 50X response, build failed",
			statusCode:    http.StatusInternalServerError,
			expectedError: "unable to build job - status: 500 body: ",
		},
	}

	for _, tt := range tests {
		s := buildHttpTestServer(tt.statusCode, "")
		jr := NewJobRunner(logger, tt.user, tt.token, s.URL)

		err := jr.RunBuild()
		if err != nil {
			fmt.Print(err)
			if tt.expectedError == "" {
				t.Errorf("unexpected error: %s", err.Error())
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %s, but got error %s", tt.expectedError, err.Error())
			}
		}
	}
}

func buildHttpTestServer(s int, b string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(s)
		_, err := w.Write([]byte(b))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Printf("%+v", r.URL)
	}))
	return ts
}
