package common

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

//HTTPToMattermost send data to mattermost
func HTTPToMattermost(body string) error {

	urlMattermost := "http://138.91.123.197:10080/hooks/sdn1z4obxfr65yt88df7aybauy"

	bodyByte := []byte(body)

	request, err := http.NewRequest("POST", urlMattermost, bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer response.Body.Close()

	fmt.Println(response.Status)

	return nil
}
