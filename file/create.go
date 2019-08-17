package file

import (
	"bufio"
	"os"
	"swgw/common"
)

//Template CloudFormatonの構造体
type Template struct {
	AWSTemplateFormatVersion string
	Description              string
	Resources                string
}

//CreateTemplate CloudFormationのテンプレートを作成する
func CreateTemplate(args []string, fileName string) {

	file, err := os.OpenFile(fileName, os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	data := Template{AWSTemplateFormatVersion: "2010-09-09"}

	common.Write(file, data)

	//引数に指定されたリソースごとにResourcesセクションを追記する
	if len(args) > 0 {
		addResources(file, args)
	}

}

//AddBlunk 先頭にスペースを2つ追加する
func AddBlunk(fileName string, resultFileName string) {
	file, err := os.Open("template.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	afterResources := false
	newFile, err := os.OpenFile(resultFileName, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	for scan.Scan() {
		line := common.ReplaceUpper(scan.Text())

		if line == `Resources: ""` {
			newFile.WriteString(line[:10] + "\n")
		} else if !afterResources {
			newFile.WriteString(line + "\n")
		} else {
			newFile.WriteString("  " + line + "\n")
		}

		if line == `Resources: ""` {
			afterResources = true
		}

	}

}

func addResources(file *os.File, args []string) {

	for _, arg := range args {
		switch arg {
		case "ecr":
			AddEcr(file)
		case "ec2":
			AddEC2(file)
		case "eip":
			AddEip(file)
		case "nat", "ng":
			AddNat(file)
		case "acl":
			AddACL(file)
		case "subnet":
			AddSubnet(file)
		case "elb", "lb":
			AddElb(file)
		case "deploy", "codedeploy", "cd":
			AddDeploy(file)
		case "ecs":
			AddCluster(file)
			AddService(file)
			AddTask(file)
		case "cluster":
			AddCluster(file)
		case "service":
			AddService(file)
		case "task", "taskdefinition", "td":
			AddTask(file)
		}
	}

}
