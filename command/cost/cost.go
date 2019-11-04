package cost

import (
	"fmt"
	"swgw/common"

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
			return getDailyCost(cmd, args, sess)
		},
	}
	cmd.AddCommand(NewCostS3Command(sess))
	return cmd
}

//AWSCostUsage サービスごとの使用量と金額を格納する構造体
type AWSCostUsage struct {
	Service string
	Cost    string
	Usage   string
}

func getDailyCost(cmd *cobra.Command, args []string, sess *session.Session) error {
	return nil
	// svc := costexplorer.New(sess)

	// var params = costexplorer.GetCostAndUsageInput{
	// 	Granularity: aws.String(costexplorer.GranularityDaily),
	// 	GroupBy:     []*costexplorer.GroupDefinition{new(costexplorer.GroupDefinition).SetKey("SERVICE").SetType("DIMENSION")},
	// 	Metrics:     aws.StringSlice([]string{"UnblendedCost", "UsageQuantity"}),
	// 	TimePeriod:  new(costexplorer.DateInterval).SetStart(time.Now().AddDate(0, 0, -1).Format(layout)).SetEnd(time.Now().Format(layout)),
	// }

	// res, err := svc.GetCostAndUsage(&params)
	// if err != nil {
	// 	return err
	// }

	// parsedRes := parse(res)
	// show(parsedRes)

	// return nil
}

//GetCostTotal 合計金額を取得する
// func GetCostTotal(sess *session.Session) {
// 	svc := costexplorer.New(sess)

// 	now := time.Now()

// 	var params = costexplorer.GetCostAndUsageInput{
// 		Granularity: aws.String(costexplorer.GranularityMonthly),
// 		GroupBy:     []*costexplorer.GroupDefinition{},
// 		Metrics:     aws.StringSlice([]string{"UnblendedCost"}),
// 		TimePeriod:  new(costexplorer.DateInterval).SetStart(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format(layout)).SetEnd(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1).Format(layout)),
// 	}
// 	res, err := svc.GetCostAndUsage(&params)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println()
// 	fmt.Printf("%s月の合計使用金額：$%s\n\n", now.AddDate(0, 0, -1).Format("01")[1:], aws.StringValue(res.ResultsByTime[0].Total["UnblendedCost"].Amount))
// }

func parse(cost *costexplorer.GetCostAndUsageOutput) []AWSCostUsage {

	resources := cost.ResultsByTime
	group := resources[0].Groups

	var costUsages []AWSCostUsage

	for _, group := range group {
		costUsages = append(costUsages,
			AWSCostUsage{
				Service: aws.StringValue(group.Keys[0]),
				Cost:    aws.StringValue(group.Metrics["UnblendedCost"].Amount),
				Usage:   aws.StringValue(group.Metrics["UsageQuantity"].Amount),
			})
	}

	return costUsages
}

func show(costUsages []AWSCostUsage) {

	for _, costUsage := range costUsages {
		fmt.Printf("■ %-40s：$%s (使用量：%s)\n", costUsage.Service, costUsage.Cost, costUsage.Usage)
	}
}
