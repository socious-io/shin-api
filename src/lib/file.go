package lib

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"shin/src/utils"
)

func DownloadFile(url string) (io.ReadCloser, error) {
	// Fetch the CSV data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return resp.Body, nil
}

func ProcessCSVFromUrl(url string, fields []string) (result []map[string]string, err error) {

	body, err := DownloadFile(url)
	if err != nil {
		fmt.Println("err", err.Error())
		return nil, err
	}
	defer body.Close()

	// Read the CSV data
	reader := csv.NewReader(body)
	var csvFields []string
	data := []map[string]string{}

	// Iterate over the rows
	c := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}

		if c == 0 {
			csvFields = record
			//validate csv fields
			for _, field := range fields {
				if !utils.ArrayContains(csvFields, field) {
					fmt.Printf("Validation failed: there's mismatch between schema and CSV fields")
					return nil, fmt.Errorf("Validation failed: there's mismatch between schema and CSV fields")
				}
			}
		} else {
			newData := map[string]string{}
			for i, value := range record {
				newData[csvFields[i]] = value
			}
			data = append(data, newData)
		}

		c++
	}

	return data, nil
}
