package file

import (
	"os"
	"swgw/common"
)

//ServiceData ServiceのCloudFormation構造体
type ServiceData struct {
	Service ServiceResources
}

//ServiceResources Resources構造体
type ServiceResources struct {
	Type       string
	Properties ServiceProperty
}

//ServiceProperty Serviceプロパティ情報
type ServiceProperty struct {
	Cluster                       string
	DeploymentConfiguration       DeployConf
	DesiredCount                  int
	HealthCheckGracePeriodSeconds int
	LaunchType                    string
	LoadBalancers                 []Lb
	NetworkConfiguration          Network
	PlacementConstraints          []Constraint
	PlacementStrategies           []Strategy
	PlatformVersion               string
	Role                          string
	SchedulingStrategy            string
	ServiceName                   string
	ServiceRegistries             []Registry
	TaskDefinition                string
}

//DeployConf サービス更新時に実行するタスクの数
type DeployConf struct {
	MaximumPercent        int
	MinimumHealthyPercent int
}

//Lb サービスと関連付けるロードバランサー
type Lb struct {
	ContainerName    string
	ContainerPort    int
	LoadBalancerName string
	TargetGroupArn   string
}

//Network サービスのネットワーク構成
type Network struct {
	AwsvpcConfiguration Awsvpn
}

//Awsvpn タスクやサービスのサブネットとセキュリティグループ
type Awsvpn struct {
	AssignPublicIP string
	SecurityGroups []string
	Subnets        []string
}

//Constraint タスク配置制約
type Constraint struct {
	Type       string
	Expression string
}

//Strategy タスク配置戦略
type Strategy struct {
	Type  string
	Filed string
}

//Registry サービスレジストリ
type Registry struct {
	ContainerName string
	ContainerPort int
	Port          int
	RegistryArn   string
}

//AddService サービスを追加する
func AddService(file *os.File) {
	var data ServiceData
	data.Service.Type = "AWS::ECS::Service"
	data.Service.Properties.LoadBalancers = append(data.Service.Properties.LoadBalancers, Lb{ContainerName: "", ContainerPort: 80, LoadBalancerName: "", TargetGroupArn: ""})
	data.Service.Properties.PlacementConstraints = append(data.Service.Properties.PlacementConstraints, Constraint{Type: "", Expression: ""})
	data.Service.Properties.PlacementStrategies = append(data.Service.Properties.PlacementStrategies, Strategy{Type: "", Filed: ""})
	data.Service.Properties.ServiceRegistries = append(data.Service.Properties.ServiceRegistries, Registry{ContainerName: "", ContainerPort: 80, Port: 80, RegistryArn: ""})
	common.Write(file, data)
}
