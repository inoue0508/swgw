package file

import (
	"os"
	"swgw/common"
)

//ACLData ACLのCloudFormation構造体
type ACLData struct {
	ACL ACLResources
}

//ACLResources Resources構造体
type ACLResources struct {
	Type       string
	Properties ACLProperty
}

//ACLProperty ACLプロパティ情報
type ACLProperty struct {
	VpcID string
	Tags  []Tag
}

//ACLEntryData ACLのルール情報
type ACLEntryData struct {
	NetworkACLEntry ACLEntryResource
}

//ACLEntryResource AclEntryのリソース情報
type ACLEntryResource struct {
	Type       string
	Properties ACLEntryProperty
}

//ACLEntryProperty Aclのルールのプロパティ情報
type ACLEntryProperty struct {
	CidrBlock     string
	Egress        bool
	Icmp          IcmpInfo
	Ipv6CidrBlock string
	NetworlACLID  string
	PortRange     PortInfo
	Protocol      int
	RuleAction    string
	RuleNumber    int
}

//IcmpInfo Icmp情報
type IcmpInfo struct {
	Code int
	Type int
}

//PortInfo ポート情報
type PortInfo struct {
	From int
	To   int
}

//SubnetACLData SubnetとAclの関連付け
type SubnetACLData struct {
	SubentACLAssociation SubnetACLResource
}

//SubnetACLResource SubnetAcl関連付けのリソース情報
type SubnetACLResource struct {
	Type       string
	Properties SubnetACLProperty
}

//SubnetACLProperty SubnetAcl関連付けのプロパティ
type SubnetACLProperty struct {
	SubnetID     string
	NetworkACLID string
}

//AddACL ACLリソースを追加する
func AddACL(file *os.File) {
	var data ACLData
	data.ACL.Type = "AWS::EC2::NetworkAcl"
	data.ACL.Properties.Tags = append(data.ACL.Properties.Tags, Tag{Key: "", Value: ""})
	common.Write(file, data)

	var entrydata ACLEntryData
	entrydata.NetworkACLEntry.Type = "AWS::EC2::NetworkAclEntry"
	common.Write(file, entrydata)

	var subnetACLData SubnetACLData
	subnetACLData.SubentACLAssociation.Type = "AWS::EC2::SubnetNetworkAclAssociation"
	common.Write(file, subnetACLData)

}
