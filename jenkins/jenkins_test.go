package jenkins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joerdav/zapray"
)

func TestJenkinsBuild(t *testing.T) {
	logger, err := zapray.NewDevelopment()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		statusCode int
		error      string
		user       string
		token      string
		url        string
		name       string
	}{
		{
			name:       "given a 201 response, build started successfully",
			statusCode: http.StatusCreated,
		},
		{
			name:       "given, a 40X response, build failed",
			statusCode: http.StatusBadRequest,
			error:      "unable to build job - status: 400 body: ",
		},
		{
			name:       "given, a 50X response, build failed",
			statusCode: http.StatusInternalServerError,
			error:      "unable to build job - status: 500 body: ",
		},
	}

	for _, tt := range tests {
		// Arrange
		url := tt.url
		if url == "" {
			// setup the JenkinsJob to use our test server
			url = buildHttpTestServer(tt.statusCode, "").URL
		}
		jj := NewJenkinsJob(logger, tt.user, tt.token, url)

		// Act
		err := jj.Build()

		// Assert
		if err != nil {
			if err.Error() != tt.error {
				t.Errorf("expected %s, but got %s", tt.error, err.Error())
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
