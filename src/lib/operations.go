package lib

import (
	"context"
	"shin/src/app/models"

	"github.com/google/uuid"
)

func RevokeCredentialOperation(ctx context.Context, credential models.Credential, doneChan chan uuid.UUID, errChan chan error) {
	if err := credential.Revoke(ctx.(context.Context)); err != nil {
		errChan <- err
		return
	}
	doneChan <- credential.ID
}
