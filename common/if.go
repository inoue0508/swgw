package common

import (
	"os"

	"github.com/go-yaml/yaml"
)

//Write 第二引数のデータを第一引数のfileに書き込む
func Write(file *os.File, data interface{}) {
	buf, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}
	file.Write(buf)
}
