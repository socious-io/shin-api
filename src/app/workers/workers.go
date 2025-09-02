package workers

import (
	"github.com/socious-io/gomq"
)

var consumers = []gomq.AddConsumerParams{
	{
		Channel:       "sociousid/event:user.delete",
		Consumer:      gomq.NewConsumer(DeleteUser),
		IsCategorized: false,
	},
	{
		Channel:       "sociousid/event:identities.sync",
		Consumer:      gomq.NewConsumer(SyncIdentities),
		IsCategorized: false,
	},
	{
		Channel:       "import",
		Consumer:      gomq.NewConsumer(ImportWorker),
		IsCategorized: true,
	},
	{
		Channel:       "operation",
		Consumer:      gomq.NewConsumer(OperationWorker),
		IsCategorized: true,
	},
}

func RegisterConsumers() {
	for _, consumer := range consumers {
		gomq.AddConsumer(consumer)
	}
}
