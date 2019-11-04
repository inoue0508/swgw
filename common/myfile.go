package common

import (
	"archive/zip"
	"io"
	"log"
	"os"
)

//Unzip unzip zip file
func Unzip(zipFile string) ([]string, error) {
	// ref
	// https://taknb2nch.hatenablog.com/entry/20140109/1389259972
	// https://gist.github.com/hnaohiro/4572580

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer zipReader.Close()

	var fileList []string
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		defer rc.Close()

		fileList = append(fileList, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(file.Name, file.Mode())
		} else {
			f, err := os.OpenFile(file.Name, os.O_WRONLY|os.O_CREATE, file.Mode())
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}
	}

	return fileList, nil
}
