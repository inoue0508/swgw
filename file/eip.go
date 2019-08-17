package file

import (
	"os"
	"swgw/common"
)

//EipData EC2のCloudFormation構造体
type EipData struct {
	Eip EipResources
}

//EipResources Resources構造体
type EipResources struct {
	Type       string
	Properties EipProperty
}

//EipProperty EIPのプロパティ
type EipProperty struct {
	Domain         string
	InstanceID     string
	PublicIpv4Pool string
}

//AddEip EIPリソースを追加する
func AddEip(file *os.File) {
	var data EipData
	data.Eip.Type = "AWS::EC2::EIP"
	common.Write(file, data)
}
