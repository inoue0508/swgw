package file

import (
	"os"
	"swgw/common"
)

//ElbData ElbのCloudFormation構造体
type ElbData struct {
	Elb ElbResources
}

//ElbResources Resources構造体
type ElbResources struct {
	Type       string
	Properties ElbProperty
}

//ElbProperty Elbプロパティ情報
type ElbProperty struct {
	IPAddressType          string
	LoadBalancerAttributes []Tag
	Name                   string
	Scheme                 string
	SecurityGroups         []string
	SubnetMappings         []SubnetMapping
	Subnets                []string
	Tags                   []Tag
	Type                   string
}

//SubnetMapping サブネットのマッピング
type SubnetMapping struct {
	SubnetID     string
	AllocationID string
}

//ListenerData ListenerのCloudFormation構造体
type ListenerData struct {
	Listener ListenerResources
}

//ListenerResources Resources構造体
type ListenerResources struct {
	Type       string
	Properties ListenerProperty
}

//ListenerProperty Listenerプロパティ情報
type ListenerProperty struct {
	Certificates    []Certificate
	DefaultAction   ActionInfo
	LoadBalancerArn string
	Port            int
	Protocol        string
	SslPolicy       string
}

//Certificate 証明書情報
type Certificate struct {
	CertificateArn string
}

//ActionInfo リスナーが実行するデフォルトアクション
type ActionInfo struct {
	AuthenticateCognitoConfig CognitoConfig
	AuthenticateOidcConfig    OidcConfig
	FixedResponseConfig       ResponseConfig
	Order                     int
	RedirectConfig            RedirectConf
	TargetGroupArn            string
	Type                      string
}

//CognitoConfig Amazon Cognitoと統合するときの情報
type CognitoConfig struct {
	AuthenticationRequestExtraParams string
	OnUnauthenticatedRequest         string
	Scope                            string
	SessionCookieName                string
	SessionTimeout                   int
	UserPoolClientID                 string
	UserPoolDomain                   string
}

//OidcConfig OIDCユーザー認証時のパラメータリクエスト
type OidcConfig struct {
	AuthenticationRequestExtraParams string
	AuthorizationEndpoint            string
	ClientID                         string
	ClientSecret                     string
	Issuer                           string
	OnUnauthenticatedRequest         string
	Scope                            string
	SessionCookieName                string
	SessionTimeout                   int
	TokenEndpoint                    string
	UserinfoEndpoint                 string
}

//ResponseConfig カスタムHTTPレスポンス
type ResponseConfig struct {
	ContentType string
	MessageBody string
	StatusCode  string
}

//RedirectConf リダイレクト情報
type RedirectConf struct {
	Host       string
	Path       string
	Protocol   string
	Query      string
	StatusCode string
}

//TargetData ElbTargetのCloudFormation構造体
type TargetData struct {
	Target TargetResources
}

//TargetResources TargetResources構造体
type TargetResources struct {
	Type       string
	Properties TargetProperty
}

//TargetProperty ElbTargetプロパティ情報
type TargetProperty struct {
	HealthCheckIntervalSeconds int
	HealthCheckPath            string
	HealthCheckPort            string
	HealthCheckProtocol        string
	HealthCheckTimeoutSeconds  int
	Matcher                    MatcherInfo
	Name                       string
	Port                       int
	Protocol                   string
	Tags                       []Tag
	TargetGroupAttributes      []Tag
	Targets                    []TargetDefinition
	TargetType                 string
	UnhealthyThresholdCount    int
	VpcID                      string
}

//MatcherInfo ヘルスチェック応答に必要なHTTPコード
type MatcherInfo struct {
	HTTPCode string
}

//TargetDefinition ターゲットグループに追加するターゲット
type TargetDefinition struct {
	AvailabilityZone string
	ID               string
	Port             string
}

//AddElb elbリソースを追加する
func AddElb(file *os.File) {
	var data ElbData
	data.Elb.Type = "AWS::ElasticLoadBalancingV2::LoadBalancer"
	data.Elb.Properties.Tags = append(data.Elb.Properties.Tags, Tag{Key: "", Value: ""})
	data.Elb.Properties.LoadBalancerAttributes = append(data.Elb.Properties.LoadBalancerAttributes, Tag{Key: "", Value: ""})
	data.Elb.Properties.SubnetMappings = append(data.Elb.Properties.SubnetMappings, SubnetMapping{SubnetID: "", AllocationID: ""})
	common.Write(file, data)

	var listenerData ListenerData
	listenerData.Listener.Type = "AWS::ElasticLoadBalancingV2::Listener"
	listenerData.Listener.Properties.Certificates = append(listenerData.Listener.Properties.Certificates, Certificate{CertificateArn: ""})
	common.Write(file, listenerData)

	var targetData TargetData
	targetData.Target.Type = "AWS::ElasticBalancingV2::TargetGroup"
	targetData.Target.Properties.Tags = append(targetData.Target.Properties.Tags, Tag{Key: "", Value: ""})
	targetData.Target.Properties.TargetGroupAttributes = append(targetData.Target.Properties.TargetGroupAttributes, Tag{Key: "", Value: ""})
	targetData.Target.Properties.Targets = append(targetData.Target.Properties.Targets, TargetDefinition{AvailabilityZone: "", ID: "", Port: ""})
	common.Write(file, targetData)

}
