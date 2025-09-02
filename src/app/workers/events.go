package workers

import (
	"context"
	"log"
	"shin/src/app/models"
)

func DeleteUser(form DeleteUserForm) error {
	ctx := context.Background()

	u, err := models.GetUser(form.User.ID)
	if err != nil {
		return err
	}

	if err := u.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func SyncIdentities(form SyncForm) error {
	ctx := context.Background()

	if err := form.User.Upsert(ctx); err != nil {
		log.Printf("SyncIdentities: Error upserting user: %v", err)
		return err
	}

	for _, o := range form.Organizations {
		if err := o.Create(ctx, form.User.ID); err != nil {
			log.Printf("SyncIdentities: Error upserting org: %v", err)
			return err
		}
	}

	return nil
}
