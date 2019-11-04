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

//SendCost struct
type SendCost struct {
	Yesterday string
	CostURL   string
}

//SendCostInfo snd yesterday's cost info to mattermost
func SendCostInfo(costData []CsvData) error {

	message := `{"text": "#### {{.Yesterday}}\n
| Service | Cost[USD]                      |
|:--------|:------------------------------|`

	for _, data := range costData {
		row := `| ` + data.ServiceName + ` | ` + data.Cost + ` |
`
		message += row
	}
	message += `[cost explore]({{.CostURL}})`

	var returnMessage bytes.Buffer
	msg, err := template.New("myTemplate").Parse(message)
	if err != nil {
		return err
	}

	replace := SendCost{
		Yesterday: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		CostURL:   "https://console.aws.amazon.com/cost-reports/home?#/custom?groupBy=Service&hasBlended=false&hasAmortized=false&excludeDiscounts=true&excludeTaggedResources=false&timeRangeOption=Custom&granularity=Daily&reportName=&reportType=CostUsage&isTemplate=true&startDate=2019-10-01&endDate=2019-11-04&filter=%5B%7B%22dimension%22:%22RecordType%22,%22values%22:%5B%22Refund%22,%22Credit%22%5D,%22include%22:false,%22children%22:null%7D%5D&forecastTimeRangeOption=None&usageAs=usageQuantity&chartStyle=Stack",
	}

	err = msg.Execute(&returnMessage, replace)
	fmt.Println(returnMessage.String())
	return nil
}
