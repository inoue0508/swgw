package file

import (
	"os"
	"swgw/common"
)

//SubnetData SubnetのCloudFormation構造体
type SubnetData struct {
	Subnet SubnetResources
}

//SubnetResources Resources構造体
type SubnetResources struct {
	Type       string
	Properties SubnetProperty
}

//SubnetProperty Subnetプロパティ情報
type SubnetProperty struct {
	AssignIpv6AddressOnCreation bool
	AvailabilityZone            string
	CidrBlock                   string
	Ipv6CidrBlock               string
	MapPublicIPOnLaunch         bool
	Tags                        []Tag
	VpcID                       string
}

//SubnetCidrData SubnetCidrのCloudFormation構造体
type SubnetCidrData struct {
	SubnetCidr SubnetCidrResources
}

//SubnetCidrResources Resources構造体
type SubnetCidrResources struct {
	Type       string
	Properties SubnetCidrProperty
}

//SubnetCidrProperty SubnetのCidrのプロパティ
type SubnetCidrProperty struct {
	Ipv6CidrBlock string
	SubnetID      string
}

//AddSubnet サブネットリソースを追加する
func AddSubnet(file *os.File) {
	var data SubnetData
	data.Subnet.Type = "AWS::EC2::Subnet"
	data.Subnet.Properties.Tags = append(data.Subnet.Properties.Tags, Tag{Key: "", Value: ""})
	common.Write(file, data)

	var cidrData SubnetCidrData
	cidrData.SubnetCidr.Type = "AWS::EC2::SubnetCidrBlock"
	common.Write(file, cidrData)
}
