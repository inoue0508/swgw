package common

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"
)

var urlMattermost = os.Getenv("MATTERMOST_URL")
var pipelineURL = os.Getenv("PIPELINE_URL")

//NotifyMM mattermostにhttpリクエストを送信する
func NotifyMM(name, status string) error {

	body, err := createJSON(name, status)
	if err != nil {
		return err
	}
	bodyByte := []byte(body)

	request, err := http.NewRequest("POST", urlMattermost, bytes.NewBuffer(bodyByte))
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

//SendS3Cost struct
type SendS3Cost struct {
	Yesterday string
	CostURL   string
}

//SendS3CostInfo snd yesterday's cost info to mattermost
func SendS3CostInfo(costData []CsvData) (string, error) {

	message := `{"text": "#### {{.Yesterday}}\n
| Service | Cost[USD]                      |
|:--------|:------------------------------|
`

	for _, data := range costData {
		row := `| ` + data.ServiceName + ` | ` + data.Cost + ` |
`
		message += row
	}
	message += `[cost explore]({{.CostURL}})"}`

	var returnMessage bytes.Buffer
	msg, err := template.New("myTemplate").Parse(message)
	if err != nil {
		return "", err
	}

	replace := SendS3Cost{
		Yesterday: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		CostURL:   "https://console.aws.amazon.com/cost-reports/home?#/dashboard",
	}

	err = msg.Execute(&returnMessage, replace)
	fmt.Println(returnMessage.String())
	return returnMessage.String(), nil
}

//AWSCostDaily2Days 2 days cost
type AWSCostDaily2Days struct {
	Service    string
	DBYCost    string
	YCost      string
	Difference string
	Sum        string
}

//AWSMonthCost 1月のコスト
type AWSMonthCost struct {
	Service string
	Cost    string
}

//CostReplace コストの置換構造体
type CostReplace struct {
	Yesterday   string
	DBYesterdey string
	CostURL     string
	Month       string
}

//CreateCostInfoBody create   http request body
func CreateCostInfoBody(awsCost []AWSCostDaily2Days) (string, error) {

	message := `{"text": "#### {{.Yesterday}}の使用状況\n
| Service | {{.Yesterday}} [$] | {{.DBYesterdey}} [$] | 差 [$] | {{.Month}}月合計 [$] |
|:--------|:-------------------|:---------------------|:-------|:---------------------|
`

	for _, cost := range awsCost {
		row := `| ` + cost.Service + ` | ` + cost.YCost + ` | ` + cost.DBYCost + ` | ` + cost.Difference + ` | ` + cost.Sum + ` |
`
		message += row
	}

	message += `[CostExploreコンソールへ]({{.CostURL}})"}`

	replace := CostReplace{
		Yesterday:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		DBYesterdey: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
		CostURL:     "https://console.aws.amazon.com/cost-reports/home?#/dashboard",
		Month:       time.Now().AddDate(0, 0, -1).Format("01"),
	}

	var returnMessage bytes.Buffer
	msg, err := template.New("myTemplate").Parse(message)
	if err != nil {
		return "", err
	}
	err = msg.Execute(&returnMessage, replace)

	return returnMessage.String(), nil

}
