package file

import (
	"os"
	"swgw/common"
)

//ApplicationData ApplicationのCloudFormation構造体
type ApplicationData struct {
	Application ApplicationResources
}

//ApplicationResources Resources構造体
type ApplicationResources struct {
	Type       string
	Properties ApplicationProperty
}

//ApplicationProperty DeployApplicationnのプロパティ
type ApplicationProperty struct {
	ApplicationName string
	ComputePlatform string
}

//DeployData CodeDeploy
type DeployData struct {
	DeployGroup DeployResource
}

//DeployResource CodeDeploy
type DeployResource struct {
	Type       string
	Properties DeployProperty
}

//DeployProperty CodeDeploy
type DeployProperty struct {
	AlarmConfiguration           AlarmConf
	ApplicationName              string
	AutoRollbackConfiguration    Rollback
	AutoScalingGroups            []string
	Deployment                   Deploy
	DeploymentConfigName         string
	DeploymentGroupName          string
	DeploymentStyle              Style
	Ec2TagFilters                []TagFilter
	Ec2TagSet                    TagSetList
	LoadBalancerInfo             LbInfo
	OnPremisesInstanceTagFilters []TagFilter
	OnPremisesInstanceTagSet     OnpreTagSet
	ServiceRoleArn               string
	TriggerConfigurations        []TriggerConf
}

//AlarmConf CloudWacthアラーム
type AlarmConf struct {
	Alarms          []Alarm
	Enabled         bool
	IgnorePollAlarm bool
}

//Alarm CloudWatchアラーム名
type Alarm struct {
	Name string
}

//Rollback ロールバック
type Rollback struct {
	Enabled bool
	Events  []string
}

//Deploy Deployのリビジョン
type Deploy struct {
	Description                   string
	IgnoreApplicationStopFailures bool
	Revision                      Revision
}

//Revision リビジョン情報
type Revision struct {
	GitHubLocation GitLoc
	RevisionType   string
	S3Location     S3Loc
}

//GitLoc gitのロケーション情報
type GitLoc struct {
	CommitID   string
	Repository string
}

//S3Loc s3情報
type S3Loc struct {
	Bucket     string
	BundleType string
	ETag       string
	Key        string
	Version    string
}

//Style デプロイ実行タイプ
type Style struct {
	DeploymentOption string
	DeploymentType   string
}

//TagFilter EC2のタグフィルター
type TagFilter struct {
	Key   string
	Type  string
	Value string
}

//TagSetList タグセットのリスト
type TagSetList struct {
	EC2TagSetListObject []TagSet
}

//TagSet タグセット情報
type TagSet struct {
	Ec2TagGroup []TagFilter
}

//LbInfo LB情報
type LbInfo struct {
	ElbInfoList         []Name
	TargetGroupInfoList []Name
}

//Name name
type Name struct {
	Name string
}

//OnpreTagSet オンプレタグセット
type OnpreTagSet struct {
	OnPremisesTagGroup []TagFilter
}

//TriggerConf トリガー情報
type TriggerConf struct {
	TriggerEvents    []string
	TriggerName      string
	TriggerTargetArn string
}

//AddDeploy CodeDeployリソースを追加する
func AddDeploy(file *os.File) {
	var data ApplicationData
	data.Application.Type = "AWS::CodeDeploy::Application"
	common.Write(file, data)

	var deployData DeployData
	deployData.DeployGroup.Type = "AWS::CodeDeploy::DeployGroup"
	tagFilter := TagFilter{Key: "", Type: "", Value: ""}
	deployData.DeployGroup.Properties.Ec2TagFilters = append(deployData.DeployGroup.Properties.Ec2TagFilters, tagFilter)
	deployData.DeployGroup.Properties.OnPremisesInstanceTagFilters = append(deployData.DeployGroup.Properties.OnPremisesInstanceTagFilters, tagFilter)
	deployData.DeployGroup.Properties.TriggerConfigurations = append(deployData.DeployGroup.Properties.TriggerConfigurations, TriggerConf{TriggerEvents: []string{"", ""}, TriggerName: "", TriggerTargetArn: ""})
	deployData.DeployGroup.Properties.AlarmConfiguration.Alarms = append(deployData.DeployGroup.Properties.AlarmConfiguration.Alarms, Alarm{Name: ""})

	common.Write(file, deployData)
}
