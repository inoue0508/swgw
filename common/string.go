package common

import (
	"strings"
)

var replaceList = map[string]string{
	"ID":   "Id",
	"API":  "Api",
	"IP":   "Ip",
	"ACL":  "Acl",
	"HTTP": "Http",
	"CPU":  "Cpu",
	"DNS":  "Dns",
}

//ToLowerAllay 配列の中身をすべて小文字にする。
func ToLowerAllay(args []string) []string {

	var resultArgs []string

	for _, arg := range args {
		resultArgs = append(resultArgs, strings.ToLower(arg))
	}

	return resultArgs

}

//ReplaceUpper golintで構造体の要素をIDとかAPIとか一部すべて大文字にしないと警告が出る
//ただCloudFormationではIDではなくIdであるため、ここで変換する
func ReplaceUpper(str string) string {
	replacedStr := str
	for key := range replaceList {
		replacedStr = strings.Replace(replacedStr, key, replaceList[key], 1)
	}

	return replacedStr
}
