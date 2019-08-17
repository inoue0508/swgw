package file

import (
	"os"
	"swgw/common"
)

//ECRData ECRのCloudFormation構造体
type ECRData struct {
	ECR Resources
}

//Resources Resources構造体
type Resources struct {
	Type       string
	Properties Property
}

//Property プロパティ構造体
type Property struct {
	LifecyclePolicy      Lifecycle
	RepositoryName       string
	RepositoryPolicyText string
}

//Lifecycle ライフサイクル構造体
type Lifecycle struct {
	LifecyclePolicyText string
	RegistryID          string
}

//AddEcr ecrのリソースを追加する
func AddEcr(file *os.File) {

	var data ECRData
	data.ECR.Type = "AWS::ECR::Repository"
	data.ECR.Properties.RepositoryPolicyText = "User action Policy. Emit OK"
	data.ECR.Properties.LifecyclePolicy.LifecyclePolicyText = "Lifecycle Policy with JSON"
	data.ECR.Properties.LifecyclePolicy.RegistryID = "AccountID. Emit OK"
	common.Write(file, data)
}
