package common

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

//CsvData csvdata
type CsvData struct {
	ServiceName string
	Cost        string
}

//ReadCsv read and get some data from svc file where time = yesterday
func ReadCsv(csvfile string) ([]CsvData, error) {

	now := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var returnData []CsvData
	file, err := os.Open(csvfile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		if strings.Contains(record[10], now) && strings.Contains(record[11], now) {
			returnData = append(returnData, CsvData{
				ServiceName: record[12],
				Cost:        record[24],
			})
		}

	}
	sort.SliceStable(returnData, func(i, j int) bool { return returnData[i].ServiceName < returnData[j].ServiceName })
	return returnData, nil

}
