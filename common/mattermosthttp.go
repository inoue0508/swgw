package common

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

var url = os.Getenv("MATTERMOST_URL")
var pipelineURL = os.Getenv("PIPELINE_URL")

//NotifyMM mattermostにhttpリクエストを送信する
func NotifyMM(name, status string) error {

	body, err := createJSON(name, status)
	if err != nil {
		return err
	}
	bodyByte := []byte(body)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println(response.Status)

	return nil

}

//Replace 置換するための構造体
type Replace struct {
	PipelineName   string
	PipelineStatus string
	URL            string
}

func createJSON(name, status string) (string, error) {

	var jsonStr string

	if status == "InProgress" {
		jsonStr = `{"text": "@channel\n#### :aws::aws:パイプラインがスタートしました :aws::aws:\n
| パイプライン | 状態                      |
|:--------|:------------------------------|
| {{.PipelineName}} | {{.PipelineStatus}} |
[CodePipelineコンソールへ]({{.URL}})"}`
	} else if status == "Succeeded" {
		jsonStr = `{"text": "@channel\n#### :aws::aws:パイプラインが成功しました :aws::aws:\n
| パイプライン | 状態                      |
|:--------|:------------------------------|
| {{.PipelineName}} | {{.PipelineStatus}} |
[CodePipelineコンソールへ]({{.URL}})"}`
	} else if status == "Failed" {
		jsonStr = `{"text": "@channel\n#### :aws::aws:パイプラインが失敗しました :aws::aws:\n
確認と修正をお願いします\n
| パイプライン | 状態                      |
|:--------|:------------------------------|
| {{.PipelineName}} | **{{.PipelineStatus}}** |
[CodePipelineコンソールへ]({{.URL}})"}`
	}

	var returnMessage bytes.Buffer
	msg, err := template.New("myTemplate").Parse(jsonStr)
	if err != nil {
		return "", err
	}

	replace := Replace{
		PipelineName:   name,
		PipelineStatus: status,
		URL:            pipelineURL,
	}

	err = msg.Execute(&returnMessage, replace)
	return returnMessage.String(), nil

}
