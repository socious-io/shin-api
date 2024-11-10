package lib

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"net/url"
	"shin/src/app/models"
	"shin/src/utils"
	"strconv"
	"time"
)

type CSVFieldType struct {
	Name string
	Type string
}

func ValidateCredentialField(field string, fieldType string, value string) (any, error) {
	switch fieldType {
	case string(models.Text):
		return value, nil
	case string(models.Number):
		convValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("Validation failed: Type of field: %s isn't %s (error: %s)", field, models.Number, err.Error())
		}
		return convValue, nil
	case string(models.Boolean):
		convValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("Validation failed: Type of field: %s isn't %s (error: %s)", field, models.Boolean, err.Error())
		}
		return convValue, nil
	case string(models.Email):
		convValue, err := mail.ParseAddress(value)
		if err != nil {
			return nil, fmt.Errorf("Validation failed: Type of field: %s isn't %s (error: %s)", field, models.Email, err.Error())
		}
		return convValue, nil
	case string(models.Url):
		convValue, err := url.ParseRequestURI(value)
		if err != nil {
			return nil, fmt.Errorf("Validation failed: Type of field: %s isn't %s (error: %s)", field, models.Url, err.Error())
		}
		return convValue, nil
	case string(models.Datetime):
		convValue, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, fmt.Errorf("Validation failed: Type of field: %s isn't %s (error: %s)", field, models.Datetime, err.Error())
		}
		return convValue, nil
	default:
		return nil, fmt.Errorf("Validation failed: Couldn't find any datatype for field: %s", field)
	}
}

func ValidateCSVFile(file multipart.File, fields map[string]string, resultChan chan bool, errChan chan error) {

	// Read the CSV data
	reader := csv.NewReader(file)
	var csvFields []CSVFieldType
	data := []map[string]any{}

	// Iterate over the rows
	c := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			errChan <- fmt.Errorf("error reading CSV: %w", err)
		}

		if c == 0 {
			//validate csv headers to be in the schema
			for key := range fields {
				if !utils.ArrayContains(record, key) {
					errChan <- fmt.Errorf("Validation failed: there's mismatch between schema and CSV fields (missing field %s in header)", key)
				}
			}
			//Maintaining CSV order for the field validation
			for _, value := range record {
				csvFields = append(csvFields, CSVFieldType{
					Name: value,
					Type: fields[value],
				})
			}
		} else {
			newData := map[string]any{}
			for i, value := range record {
				validValid, err := ValidateCredentialField(csvFields[i].Name, csvFields[i].Type, value)
				if err != nil {
					errChan <- err
					return
				}
				newData[csvFields[i].Name] = validValid
			}
			data = append(data, newData)
		}

		c++
	}

	resultChan <- true
}
