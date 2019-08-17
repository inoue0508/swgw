package file

import (
	"os"
	"swgw/common"
)

//EC2Data EC2のCloudFormation構造体
type EC2Data struct {
	EC2 EC2Resources
}

//EC2Resources Resources構造体
type EC2Resources struct {
	Type       string
	Properties EC2Property
}

//EC2Property プロパティ構造体
type EC2Property struct {
	Affinity                          string
	AvailabilityZone                  string
	BlockDeviceMappings               []BlockDeviceMap
	CreditSpecification               CreditSpec
	DisableAPITermination             bool
	EbsOptimized                      bool
	ElasticGpuSpecifications          []GpuSpec
	ElasticInferenceAccelerator       []Accelerator
	HostID                            string
	IamInstanceProfile                string
	ImageID                           string
	InstanceInitiatedShutdownBehavior string
	InstanceType                      string
	Ipv6AddressCount                  int
	Ipv6Addresses                     []Ipv6
	KernelID                          string
	KeyName                           string
	LaunchTemplate                    LaunchTemplateSpec
	LicenseSpecifications             []License
	Monitoring                        bool
	NetworkInterfaces                 NetworkInterface
	PlacementGroupName                string
	PrivateIPAddress                  string
	RamdiskID                         string
	SecurityGroupIDs                  []string
	SecurityGroups                    []string
	SourceDestChek                    bool
	SsmAssociations                   []AssociationParam
	SubnetID                          string
	Tags                              []Tag
	Tenancy                           string
	UserData                          string
	Volumes                           []Mount
	AdditionalInfo                    string
}

//BlockDeviceMap ブロックデバイスの情報をマッピングする
type BlockDeviceMap struct {
	DeviceName  string
	Ebs         BlockDevice
	NoDevice    bool
	VirtualName string
}

//BlockDevice ブロックデバイス情報
type BlockDevice struct {
	DeleteOnTermination bool
	Encrypted           bool
	Iops                int
	SnapshotID          string
	VolumeSize          string
	VolumeType          string
}

//CreditSpec CPU情報
type CreditSpec struct {
	CPUCredits string
}

//GpuSpec Elastic GPUのスペック
type GpuSpec struct {
	Type string
}

//Accelerator Accelerator情報
type Accelerator struct {
	Type string
}

//Ipv6 iipv6アドレス情報
type Ipv6 struct {
	Ipv6Address string
}

//LaunchTemplateSpec 起動時テンプレート情報
type LaunchTemplateSpec struct {
	LaunchTemplateID   string
	LaunchTemplateName string
	Version            string
}

//License ライセンス
type License struct {
	LicenseConfigurationArn string
}

//NetworkInterface ネットワークインターフェースの情報
type NetworkInterface struct {
	AssociatePublicIPAddress       bool
	DeleteOnTermination            bool
	Description                    string
	DeviceIndex                    string
	GroupSet                       []string
	NetworkInterfaceID             string
	Ipv6AddressCount               int
	Ipv6Addresses                  Ipv6
	PrivateIPAddress               string
	PrivateIPAddresses             PrivateIP
	SecondaryPrivateIPAddressCount int
	SubnetID                       string
}

//PrivateIP プライベートIPアドレス情報
type PrivateIP struct {
	PrivateIPAddress string
	Primary          bool
}

//AssociationParam SSMとパラメータを関連付ける
type AssociationParam struct {
	AssociationParameters []Parameter
	DocumentName          string
}

//Parameter SSMのパラメータ
type Parameter struct {
	Key   string
	Value []string
}

//Tag タグ情報
type Tag struct {
	Key   string
	Value string
}

//Mount マウント情報
type Mount struct {
	Device   string
	VolumeID string
}

//AddEC2 EC2のリソースを追記する
func AddEC2(file *os.File) {
	var data EC2Data
	data.EC2.Type = "AWS::EC2::Instance"

	tag := Tag{Key: "aa", Value: "bb"}
	data.EC2.Properties.Tags = append(data.EC2.Properties.Tags, tag)
	common.Write(file, data)

}
