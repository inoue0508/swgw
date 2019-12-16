package cost

import (
	"fmt"
	"os"
	"strconv"
	"swgw/common"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/spf13/cobra"
)

const layout = "2006-01-02"

//NewCostCommand costコマンドを作成する
func NewCostCommand() *cobra.Command {
	sess := common.GetAwsSession()
	cmd := &cobra.Command{
		Use:   "cost",
		Short: "not used now",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCost(cmd, args, sess)
		},
	}
	//don't use this command
	//cmd.AddCommand(NewCostS3Command(sess))
	return cmd
}

func executeCost(cmd *cobra.Command, args []string, sess *session.Session) error {

	dailyCost, err := getDailyCost(cmd, sess)
	if err != nil {
		return err
	}
	monthCost, err := getCostTotal(sess)
	if err != nil {
		return err
	}
	resultCost := combine(dailyCost, monthCost)
	body, err := common.CreateCostInfoBody(resultCost)

	url := os.Getenv("MATTERMOST_URL")
	if url == "" {
		fmt.Println("mattermostのincoming web hookが設定されていません。環境変数に追加してください")
	}

	err = common.HTTPToMattermost(body, url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

//AWSCostUsage サービスごとの使用量と金額を格納する構造体
type AWSCostUsage struct {
	Service string
	Cost    string
	Usage   string
}

//AWSCostDaily2Days 2 days cost
type AWSCostDaily2Days struct {
	Service    string
	DBYCost    string
	YCost      string
	Difference string
}

func getDailyCost(cmd *cobra.Command, sess *session.Session) ([]common.AWSCostDaily2Days, error) {

	svc := costexplorer.New(sess)

	var params = costexplorer.GetCostAndUsageInput{
		Granularity: aws.String(costexplorer.GranularityDaily),
		GroupBy:     []*costexplorer.GroupDefinition{new(costexplorer.GroupDefinition).SetKey("SERVICE").SetType("DIMENSION")},
		Metrics:     aws.StringSlice([]string{"UnblendedCost"}),
		TimePeriod:  new(costexplorer.DateInterval).SetStart(time.Now().AddDate(0, 0, -2).Format(layout)).SetEnd(time.Now().Format(layout)),
		Filter:      new(costexplorer.Expression).SetNot(new(costexplorer.Expression).SetDimensions(new(costexplorer.DimensionValues).SetKey("RECORD_TYPE").SetValues(aws.StringSlice([]string{"Credit", "Refund"})))),
	}

	res, err := svc.GetCostAndUsage(&params)
	if err != nil {
		return nil, err
	}

	dbycost, ycost := parse(res)
	diffcost := createDifference(dbycost, ycost)

	return diffcost, nil
}

//getCostTotal 合計金額を取得する
func getCostTotal(sess *session.Session) ([]common.AWSMonthCost, error) {
	svc := costexplorer.New(sess)

	now := time.Now()

	var params = costexplorer.GetCostAndUsageInput{
		Granularity: aws.String(costexplorer.GranularityMonthly),
		GroupBy:     []*costexplorer.GroupDefinition{new(costexplorer.GroupDefinition).SetKey("SERVICE").SetType("DIMENSION")},
		Metrics:     aws.StringSlice([]string{"UnblendedCost"}),
		TimePeriod:  new(costexplorer.DateInterval).SetStart(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format(layout)).SetEnd(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1).Format(layout)),
		Filter:      new(costexplorer.Expression).SetNot(new(costexplorer.Expression).SetDimensions(new(costexplorer.DimensionValues).SetKey("RECORD_TYPE").SetValues(aws.StringSlice([]string{"Credit", "Refund"})))),
	}
	res, err := svc.GetCostAndUsage(&params)
	if err != nil {
		return nil, err
	}

	var monthCost []common.AWSMonthCost
	sum := 0
	for _, cost := range res.ResultsByTime[0].Groups {
		monthCost = append(monthCost,
			common.AWSMonthCost{
				Service: aws.StringValue(cost.Keys[0]),
				Cost:    aws.StringValue(cost.Metrics["UnblendedCost"].Amount),
			})
		cost, _ := strconv.Atoi(aws.StringValue(cost.Metrics["UnblendedCost"].Amount))
		sum += cost
	}
	monthCost = append(monthCost,
		common.AWSMonthCost{
			Service: "合計",
			Cost:    strconv.Itoa(sum),
		})

	return monthCost, nil
}

func parse(cost *costexplorer.GetCostAndUsageOutput) ([]AWSCostUsage, []AWSCostUsage) {

	resources := cost.ResultsByTime
	groupdby := resources[0].Groups
	groupy := resources[1].Groups

	var dbyesterday []AWSCostUsage
	var yesterday []AWSCostUsage
	sumdby := 0
	for _, group := range groupdby {
		dbyesterday = append(dbyesterday,
			AWSCostUsage{
				Service: aws.StringValue(group.Keys[0]),
				Cost:    aws.StringValue(group.Metrics["UnblendedCost"].Amount),
			})
		cost, _ := strconv.Atoi(aws.StringValue(group.Metrics["UnblendedCost"].Amount))
		sumdby += cost
	}
	dbyesterday = append(dbyesterday,
		AWSCostUsage{
			Service: "合計",
			Cost:    strconv.Itoa(sumdby),
		})

	sumy := 0
	for _, group := range groupy {
		yesterday = append(yesterday,
			AWSCostUsage{
				Service: aws.StringValue(group.Keys[0]),
				Cost:    aws.StringValue(group.Metrics["UnblendedCost"].Amount),
			})
		cost, _ := strconv.Atoi(aws.StringValue(group.Metrics["UnblendedCost"].Amount))
		sumy += cost
	}
	yesterday = append(yesterday,
		AWSCostUsage{
			Service: "合計",
			Cost:    strconv.Itoa(sumy),
		})

	return dbyesterday, yesterday
}

func createDifference(dbyCost, yCost []AWSCostUsage) []common.AWSCostDaily2Days {

	var awscost []common.AWSCostDaily2Days

	for _, dby := range dbyCost {
		var diff string
		for _, y := range yCost {
			if dby.Service == y.Service {
				dbycost, _ := strconv.Atoi(dby.Cost)
				ycost, _ := strconv.Atoi(y.Cost)
				diff = strconv.Itoa(ycost - dbycost)
				break
			}
		}

		if diff == "" {
			if dby.Cost != "0" {
				diff = "-" + dby.Cost
			} else {
				diff = dby.Cost
			}

		}
		awscost = append(awscost,
			common.AWSCostDaily2Days{
				Service:    dby.Service,
				DBYCost:    dby.Cost,
				YCost:      "-",
				Difference: diff,
			})
	}

	for _, y := range yCost {
		isin := false
		count := 0
		for _, sum := range awscost {
			if y.Service == sum.Service {
				isin = true
				awscost[count].YCost = y.Cost
				break
			}
			count++
		}
		if isin == false {
			awscost = append(awscost,
				common.AWSCostDaily2Days{
					Service:    y.Service,
					DBYCost:    "-",
					YCost:      y.Cost,
					Difference: y.Cost,
				})
		}
	}

	return awscost

}

func combine(daily []common.AWSCostDaily2Days, monthly []common.AWSMonthCost) []common.AWSCostDaily2Days {

	count := 0
	for _, day := range daily {
		for _, mon := range monthly {
			if day.Service == mon.Service {
				daily[count].Sum = mon.Cost
				break
			}
		}
		count++
	}

	return daily

}
