package validate

//IsAwsResources 引数がawsリソース名かをチェックする
func IsAwsResources(args, aws []string) ([]string, int) {
	var errorList []string

	for _, arg := range args {
		isResource := false
		for _, resource := range aws {
			if arg == resource {
				isResource = true
				break
			}
		}
		if isResource == false {
			errorList = append(errorList, arg)
		}
	}

	return errorList, len(errorList)

}
