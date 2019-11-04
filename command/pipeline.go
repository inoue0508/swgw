package command

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"swgw/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/cobra"
)

//InProgress pipelineが進行中の状態を表す
const InProgress = "InProgress"

var redisURL = os.Getenv("REDIS_URL")
var httpProxy = os.Getenv("PROXY")

//NewPLCommand cfコマンドを作成する
func NewPLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pl PipelineName [,PipelineName...]",
		Short: "get pipeline status and notify mattermost",
		RunE: func(cmd *cobra.Command, args []string) error {
			//sessionの作成。~/.aws/にある認証情報を使用する。
			//sess := session.Must(session.NewSessionWithOptions(session.Options{
			//	SharedConfigState: session.SharedConfigEnable,
			//}))
			httpClient := http.Client{
				Transport: &http.Transport{
					Proxy: func(*http.Request) (*url.URL, error) {
						return url.Parse(httpProxy)
					},
				},
			}
			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config: aws.Config{
					HTTPClient: &httpClient,
				},
			}))

			return notifyPipelineStatus(cmd, args, sess)
		},
	}
	return cmd
}

func notifyPipelineStatus(cmd *cobra.Command, args []string, sess *session.Session) error {

	status, err := getPipelineStatus(sess, args)
	if err != nil {
		return err
	}
	fmt.Println(status)

	red, err := redis.Dial("tcp", redisURL)
	if err != nil {
		return err
	}
	defer red.Close()

	for key := range status {
		isKey, err := redis.Bool(red.Do("EXISTS", key))
		if err != nil {
			return err
		}
		if status[key] != InProgress && isKey == false {
			//common.NotifyMM(key, status[key])
			continue
		} else if status[key] == InProgress && isKey == false {
			//ここでパイプラインがスタートになった
			_, err := red.Do("SET", key, status[key])
			if err != nil {
				return err
			}
			//スタートしたリクエストを送る
			common.NotifyMM(key, status[key])
		} else if status[key] != InProgress && isKey == true {
			//パイプラインが終了した
			_, err := red.Do("DEL", key)
			if err != nil {
				return err
			}
			//終了したリクエストを贈る
			common.NotifyMM(key, status[key])
		}

	}

	return nil

}

func getPipelineStatus(sess *session.Session, args []string) (map[string]string, error) {

	pipe := codepipeline.New(sess)

	result := make(map[string]string)
	for _, pipeName := range args {
		var param = codepipeline.GetPipelineStateInput{
			Name: aws.String(pipeName),
		}

		pipelineInfo, err := pipe.GetPipelineState(&param)
		if err != nil {
			return nil, err
		}
		result[pipeName] = aws.StringValue(pipelineInfo.StageStates[len(pipelineInfo.StageStates)-1].ActionStates[0].LatestExecution.Status)
		//fmt.Println(pipelineInfo)
		fmt.Println(aws.StringValue(pipelineInfo.StageStates[len(pipelineInfo.StageStates)-1].ActionStates[0].ActionName))
	}

	return result, nil

}
