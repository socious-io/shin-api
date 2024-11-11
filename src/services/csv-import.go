package services

import (
	"fmt"
	"shin/src/app/models"
	"shin/src/lib"
	"shin/src/utils"
)

var CSVImportChannel = CategorizeChannel("csv_import")
var CSVImportWorkersCount = 4

type CSVImportConfig struct {
	Record  map[string]any
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

func InitiateCSVImport(results []map[string]any, meta map[string]any) {

	for _, result := range results {
		SendCSVImport(CSVImportConfig{
			Record:  result,
			DocType: models.CSVDocTypeCredentials,
			Meta:    meta,
		})
	}
}
