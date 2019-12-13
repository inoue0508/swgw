package cost

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"swgw/common"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
)

const s3bucket = "billusage-swgw-poc"
const dateformat = "200601"
const tempzipfile = "tmp.csv.zip"
const tmpfile = "tmp.csv"

//NewCostS3Command get cost information from s3
func NewCostS3Command(sess *session.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "calclate day's cost",
		RunE: func(cmd *cobra.Command, args []string) error {
			return calcCost(sess)
		},
	}
	return cmd
}

func calcCost(sess *session.Session) error {
	svc := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)
	costFile, err := getCostFromS3(svc)
	if err != nil {
		return err
	}
	fmt.Println(costFile)
	err = downloadFromS3(downloader, costFile)
	if err != nil {
		return err
	}

	fileList, err := common.Unzip(tempzipfile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	csvData, err := common.ReadCsv(fileList[0])
	if err != nil {
		log.Fatal(err)
		return err
	}

	sumdata, err := sum(csvData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(sumdata)
	body, err := common.SendCostInfo(sumdata)
	if err != nil {
		log.Fatal(err)
		return err
	}

	common.HTTPToMattermost(body)

	return nil
}

func getCostFromS3(svc *s3.S3) (string, error) {

	//
	prefix := "cost/swgw-poc-report/" + getPrefixTime()
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3bucket),
		Prefix: aws.String(prefix),
	}

	listObjects, err := svc.ListObjectsV2(params)
	if err != nil {
		return "", err
	}

	var costFile string
	tmpDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, object := range listObjects.Contents {
		if strings.Contains(aws.StringValue(object.Key), ".csv.zip") {
			if aws.TimeValue(object.LastModified).After(tmpDate) {
				costFile = aws.StringValue(object.Key)
				tmpDate = aws.TimeValue(object.LastModified)
			}
		}

	}

	return costFile, nil
}

func getPrefixTime() string {

	now := time.Now().AddDate(0, 0, -1).Format(dateformat)
	netxMonth := time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0).Format(dateformat)
	return now + "01-" + netxMonth + "01/"
}

//not used now
func readS3Data(svc *s3.S3, key string) {
	key = "latest/marged_cost.csv"
	sql := "SELECT * FROM S3Object LIMIT 5"
	params := &s3.SelectObjectContentInput{
		Bucket:         aws.String(s3bucket),
		Key:            aws.String(key),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		Expression:     aws.String(sql),
		InputSerialization: &s3.InputSerialization{
			CompressionType: aws.String("NONE"),
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
				FieldDelimiter: aws.String(","),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{},
		},
	}

	res, err := svc.SelectObjectContent(params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}

func downloadFromS3(downloader *s3manager.Downloader, key string) error {
	fs, err := os.Create(tempzipfile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer fs.Close()

	param := &s3.GetObjectInput{
		Bucket: aws.String(s3bucket),
		Key:    aws.String(key),
	}
	if _, err := downloader.Download(fs, param); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func sum(csvdata []common.CsvData) ([]common.CsvData, error) {

	//add cost
	svc := ""
	var cost float64 = 0.0
	var sumData []common.CsvData
	count := 1
	for _, data := range csvdata {
		if data.ServiceName == svc {
			costfloat, _ := strconv.ParseFloat(data.Cost, 64)
			cost += costfloat
		} else if data.ServiceName != svc && svc != "" {
			sumData = append(sumData, common.CsvData{
				ServiceName: svc,
				Cost:        fmt.Sprint(cost),
			})
			svc = data.ServiceName
			costfloat, _ := strconv.ParseFloat(data.Cost, 64)
			cost = costfloat
		} else if svc == "" {
			svc = data.ServiceName
			costfloat, _ := strconv.ParseFloat(data.Cost, 64)
			cost = costfloat
		}

		if count == len(csvdata) {
			sumData = append(sumData, common.CsvData{
				ServiceName: data.ServiceName,
				Cost:        fmt.Sprint(cost),
			})
		}

		count++

	}

	return sumData, nil

}
