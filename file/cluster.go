package file

import (
	"os"
	"swgw/common"
)

//ClusterData ClusterのCloudFormation構造体
type ClusterData struct {
	Cluster ClusterResources
}

//ClusterResources Resources構造体
type ClusterResources struct {
	Type       string
	Properties ClusterProperty
}

//ClusterProperty Clusterプロパティ情報
type ClusterProperty struct {
	ClusterName string
}

//AddCluster クラスターリソースを追加する
func AddCluster(file *os.File) {
	var data ClusterData
	data.Cluster.Type = "AWS::ECS::Cluster"
	common.Write(file, data)
}
