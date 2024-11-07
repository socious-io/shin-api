package services

import (
	"fmt"
	"shin/src/app/models"
	"shin/src/lib"
	"shin/src/utils"

	"github.com/google/uuid"
)

var CSVImportChannel = CategorizeChannel("csv_import")
var CSVImportWorkersCount = 4

type CSVImportConfig struct {
	RecordID uuid.UUID
	DocType  models.CSVImportDocType
	Meta     map[string]any
	Fields   []map[string]string
}

func SendCSVImport(importConfig CSVImportConfig) {
	Mq.sendJson(CSVImportChannel, importConfig)
}

func CSVImportWorker(message interface{}) {
	importConfig := new(CSVImportConfig)
	utils.Copy(message, importConfig)

	var (
		docType  = importConfig.DocType
		recordId = importConfig.RecordID
		meta     = importConfig.Meta
	)

	if docType == models.CSVDocTypeCredentials {
		//Sending email with template
		err := lib.ImportCredentials(recordId, meta)
		if err != nil {
			fmt.Println("Couldn't Import CSVs to Database, Error: ", err.Error())
		}
	}

}
