package jenkins

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/joerdav/zapray"
)

type JenkinsJob struct {
	Log                  *zapray.Logger
	JenkinsApiToken      string
	JenkinsApiUser       string
	JenkinsBuildEndpoint string
}

func NewJenkinsJob(log *zapray.Logger, user string, token string, endpoint string) (j JenkinsJob) {
	j.JenkinsApiUser = user
	j.JenkinsApiToken = token
	j.JenkinsBuildEndpoint = endpoint
	j.Log = log
	return
}

func (j JenkinsJob) Build() (err error) {
	j.Log.Info("Building Jenkins job")
	endpoint, err := j.buildEndpoint()
	if err != nil {
		return
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(j.JenkinsApiUser, j.JenkinsApiToken)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		// body is usually empty on success, but populated on a failed request.
		body, _ := io.ReadAll(res.Body)
		err = errors.New(fmt.Sprintf("unable to build job - status: %v body: %v", res.StatusCode, string(body)))
	}
	return
}

func (j JenkinsJob) buildEndpoint() (string, error) {
	endpoint, err := url.Parse(j.JenkinsBuildEndpoint + "/build")
	if err != nil {
		return "", err
	}
	query := endpoint.Query()
	query.Add("token", j.JenkinsApiToken)
	endpoint.RawQuery = query.Encode()
	return endpoint.String(), nil
}
