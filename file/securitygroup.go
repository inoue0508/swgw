package file

import (
	"os"
	"swgw/common"
)

//SgData SgのCloudFormation構造体
type SgData struct {
	Sg SgResources
}

//SgResources Resources構造体
type SgResources struct {
	Type       string
	Properties SgProperty
}

//SgProperty Sgプロパティ情報
type SgProperty struct {
	GroupName            string
	GroupDescription     string
	SecurityGroupEgress  []SecurityRule
	SecurityGroupIngress []SecurityRule
	Tags                 []Tag
	VpcID                string
}

//SecurityRule セキュリティルール
type SecurityRule struct {
	CidrIP                     string
	CidrIpv6                   string
	Description                string
	FromPort                   int
	IPProtocol                 string
	SourceSecurityGroupID      string
	SourceSecurityGroupName    string
	SourceSecurityGroupOwnerID string
	ToPort                     int
}

//AddSg セキュリティグループリソースを追加する
func AddSg(file *os.File) {
	var data SgData
	data.Sg.Type = "AWS::EC2::SecurityGroup"
	data.Sg.Properties.Tags = append(data.Sg.Properties.Tags, Tag{Key: "", Value: ""})
	rule := SecurityRule{
		CidrIP:                     "",
		CidrIpv6:                   "",
		Description:                "",
		FromPort:                   0,
		IPProtocol:                 "",
		SourceSecurityGroupID:      "",
		SourceSecurityGroupName:    "",
		SourceSecurityGroupOwnerID: "",
		ToPort:                     0,
	}
	data.Sg.Properties.SecurityGroupEgress = append(data.Sg.Properties.SecurityGroupEgress, rule)
	data.Sg.Properties.SecurityGroupIngress = append(data.Sg.Properties.SecurityGroupIngress, rule)
	common.Write(file, data)
}
