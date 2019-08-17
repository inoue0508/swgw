package file

import (
	"os"
	"swgw/common"
)

//NatData NatのCloudFormation構造体
type NatData struct {
	Nat NatResources
}

//NatResources Resources構造体
type NatResources struct {
	Type       string
	Properties NatProperty
}

//NatProperty Natプロパティ情報
type NatProperty struct {
	AllocationID string
	SubnetID     string
	Tags         []Tag
}

//AddNat Natリソースを追加する
func AddNat(file *os.File) {
	var data NatData
	data.Nat.Type = "AWS::EC2::NatGateway"
	tag := Tag{Key: "", Value: ""}
	data.Nat.Properties.Tags = append(data.Nat.Properties.Tags, tag)
	common.Write(file, data)
}
