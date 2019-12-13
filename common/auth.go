package common

import (
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var proxy = os.Getenv("HTTP_PROXY")

//GetAwsSession get session of aws
func GetAwsSession() *session.Session {

	var sess *session.Session

	if proxy == "" {
		//no proxy
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		//proxy
		httpClient := http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return url.Parse(proxy)
				},
			},
		}

		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				HTTPClient: &httpClient,
			},
		}))
	}

	return sess
}
