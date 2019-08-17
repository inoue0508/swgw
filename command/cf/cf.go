package cf

import (
	"os"
	"swgw/common"
	"swgw/file"

	"github.com/spf13/cobra"
)

var resources = []string{
	"ec2",
	"ecr",
	"elb",
	"subnet",
	"vpc",
	"nat",
	"route",
	"rt",
	"sg",
	"lb",
	"s3",
	"rds",
}

const fileName = "template.yml"

//NewCFCommand cfコマンドを作成する
func NewCFCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cf [AWSRESOURCES,...]",
		Short: "create cloudformation template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTemplateFile(cmd, args)
		},
	}
	return cmd
}

func createTemplateFile(cmd *cobra.Command, args []string) error {

	//配列内の全てを小文字にする
	argsLower := common.ToLowerAllay(args)

	//awsリソース以外の文字列が含まれているかをチェックする
	//errorList, count := validate.IsAwsResources(argsLower, resources)
	//if count != 0 {
	//	return fmt.Errorf("以下の引数はAWSリソースとして認識されません\n%s", errorList)
	//}

	f, err := cmd.Flags().GetString("file")
	if err != nil {
		panic(err)
	}

	file.CreateTemplate(argsLower, fileName)
	file.AddBlunk(fileName, f)

	//最後にtmpファイルを削除する
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
	return nil
}
