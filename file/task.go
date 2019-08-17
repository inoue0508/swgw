package file

import (
	"os"
	"swgw/common"
)

//TaskData TaskのCloudFormation構造体
type TaskData struct {
	Task TaskResources
}

//TaskResources Resources構造体
type TaskResources struct {
	Type       string
	Properties TaskProperty
}

//TaskProperty Taskプロパティ情報
type TaskProperty struct {
	ContainerDefinitions    []Definition
	CPU                     string
	ExecutionRoleArn        string
	Family                  string
	Memory                  string
	NetworkMode             string
	PlacementConstraints    []Constraint
	RequiresCompatibilities []string
	TaskRoleArn             string
	Volumes                 []VolumeInfo
}

//Definition タスク定義
type Definition struct {
	Command                []string
	CPU                    int
	DisableNetworking      bool
	DNSSearchDomains       []string
	DNSServers             []string
	DockerLabels           string
	DockerSecurityOptions  []string
	EntryPoint             []string
	Environment            []KVPair
	Essential              bool
	ExtraHosts             []Host
	HealthCheck            Health
	Hostname               string
	Image                  string
	Links                  []string
	LinuxParameters        Linux
	LogConfiguration       Log
	Memory                 int
	MemoryReservation      int
	MountPoints            []MountInfo
	Name                   string
	PortMappings           []Port
	Privileged             bool
	ReadonlyRootFilesystem bool
	RepositoryCredentials  Cred
	Ulimits                []Ulimit
	VolumesFrom            []Volume
	WorkingDirectory       string
}

//KVPair KeyValuePair
type KVPair struct {
	Name  string
	Value string
}

//Host Host
type Host struct {
	Hostname  string
	IPAddress string
}

//Health ヘルスチェック
type Health struct {
	Command     []string
	Interval    int
	Retries     int
	StartPeriod int
	Timeout     int
}

//Linux Linux固有オプション
type Linux struct {
	Capabilities       Cap
	Devices            []Device
	InitProcessEnabled bool
	SharedMemorySize   int
	Tmpfs              []Tmp
}

//Cap ケーパビリティ
type Cap struct {
	Add  []string
	Drop []string
}

//Device デバイス
type Device struct {
	ContainerPath string
	HostPath      string
	Permissions   []string
}

//Tmp tmpfs
type Tmp struct {
	ContainerPath string
	MountOptions  []string
	Size          int
}

//Log ログ情報
type Log struct {
	LogDriver string
	Options   []string
}

//MountInfo マウント情報
type MountInfo struct {
	ContainerPath string
	SourceVolume  string
	ReadOnly      bool
}

//Port ポートマッピング
type Port struct {
	ContainerPort int
	HostPort      int
	Protocol      string
}

//Cred クレデンシャル情報
type Cred struct {
	CredentialsParameter string
}

//Ulimit Ulimit
type Ulimit struct {
	HardLimit int
	Name      string
	SoftLimit int
}

//Volume Volume
type Volume struct {
	SourceContainer string
	ReadOnly        bool
}

//VolumeInfo ボリューム
type VolumeInfo struct {
	DockerVolumeConfiguration VolumeConf
	Host                      VolumeHost
	Name                      string
}

//VolumeConf VolumeConf
type VolumeConf struct {
	Autoprovision bool
	Driver        string
	DriverOpts    []string
	Labels        []string
	Scope         string
}

//VolumeHost VolumeHost
type VolumeHost struct {
	SourcePath string
}

//AddTask タスク定義を追加する
func AddTask(file *os.File) {
	var data TaskData
	data.Task.Type = "AWS::ECS::TaskDefinition"
	data.Task.Properties.PlacementConstraints = append(data.Task.Properties.PlacementConstraints, Constraint{Type: "", Expression: ""})
	data.Task.Properties.Volumes = append(data.Task.Properties.Volumes, VolumeInfo{DockerVolumeConfiguration: VolumeConf{}, Host: VolumeHost{}, Name: ""})

	var linux Linux
	linux.Devices = append(linux.Devices, Device{ContainerPath: "", HostPath: "", Permissions: []string{"", ""}})
	linux.Tmpfs = append(linux.Tmpfs, Tmp{ContainerPath: "", MountOptions: []string{"", ""}, Size: 0})

	var task Definition
	task.Environment = append(task.Environment, KVPair{Name: "", Value: ""})
	task.ExtraHosts = append(task.ExtraHosts, Host{Hostname: "", IPAddress: ""})
	task.MountPoints = append(task.MountPoints, MountInfo{ContainerPath: "", SourceVolume: "", ReadOnly: false})
	task.PortMappings = append(task.PortMappings, Port{ContainerPort: 80, HostPort: 80, Protocol: "tcp"})
	task.Ulimits = append(task.Ulimits, Ulimit{HardLimit: 0, Name: "", SoftLimit: 0})
	task.VolumesFrom = append(task.VolumesFrom, Volume{SourceContainer: "", ReadOnly: false})
	task.LinuxParameters = linux

	data.Task.Properties.ContainerDefinitions = append(data.Task.Properties.ContainerDefinitions, task)

	common.Write(file, data)
}
