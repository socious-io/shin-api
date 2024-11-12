package services

import (
	"context"
	"fmt"
	"shin/src/app/models"
	"shin/src/utils"
)

var OperationChannel = CategorizeChannel("operation")

type OperationConfig struct {
	Trigger models.OperationServiceTrigger
	Entity  interface{}
}

func SendOperation(importConfig OperationConfig) {
	Mq.sendJson(OperationChannel, importConfig)
}

func OperationWorker(message interface{}) {
	config := new(OperationConfig)
	utils.Copy(message, config)

	var (
		trigger = config.Trigger
		entity  = config.Entity
	)

	if trigger == models.OperationCredentialRevoke {
		credentialRevoke(entity)
		return
	}

}

/*
	Operation Definitions
*/

func credentialRevoke(entity interface{}) {

	ctx := context.Background()
	credential := new(models.Credential)
	utils.Copy(entity, credential)

	if err := credential.Revoke(ctx.(context.Context)); err != nil {
		fmt.Println("Couldn't Revoke Credential, Error: ", err.Error())
		return
	}
	return

}
