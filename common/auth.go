package common

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

//GetAwsSession get session of aws
func GetAwsSession() *session.Session {
	//no proxy
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// //proxy
	// httpClient := http.Client{
	// 	Transport: &http.Transport{
	// 		Proxy: func(*http.Request) (*url.URL, error) {
	// 			return url.Parse("http://**:8080")
	// 		},
	// 	},
	// }

	// sess = session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// 	Config: aws.Config{
	// 		HTTPClient: &httpClient,
	// 	},
	// }))

	return sess
}
