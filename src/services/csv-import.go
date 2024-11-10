package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"shin/src/app/models"
	"shin/src/lib"
	"shin/src/utils"
)

var CSVImportChannel = CategorizeChannel("csv_import")
var CSVImportWorkersCount = 4

type CSVImportConfig struct {
	Record  map[string]string
	DocType models.CSVImportDocType
	Meta    map[string]any
}

func SendCSVImport(importConfig CSVImportConfig) {
	Mq.sendJson(CSVImportChannel, importConfig)
}

func CSVImportWorker(message interface{}) {
	importConfig := new(CSVImportConfig)
	utils.Copy(message, importConfig)

	var (
		docType = importConfig.DocType
		record  = importConfig.Record
		meta    = importConfig.Meta
	)

	if docType == models.CSVDocTypeCredentials {
		err := lib.ImportCredentials(record, meta)
		if err != nil {
			fmt.Println("Couldn't Import CSVs to Database, Error: ", err.Error())
		}
	}

}

/*
Initiate Import
*/

func InitiateCSVImport(file multipart.File, meta map[string]any) {

	// Read the CSV data
	file.Seek(0, 0)
	reader := csv.NewReader(file)
	var csvFields []string

	// Iterate over the rows
	c := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}

		if c == 0 {
			csvFields = record
		} else {
			//send to worker one by one
			newData := map[string]string{}
			for i, value := range record {
				newData[csvFields[i]] = value
			}

			SendCSVImport(CSVImportConfig{
				Record:  newData,
				DocType: models.CSVDocTypeCredentials,
				Meta:    meta,
			})
		}

		c++
	}
}
