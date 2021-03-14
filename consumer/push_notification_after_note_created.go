package consumer

import (
	"context"
	"fooddlv/appctx"
	"fooddlv/pubsub"
	"log"
)

// Job and Group are not used.
// Define consumer is quite hard.

//func RunPushNotificationAfterCreateNote(appCtx appctx.AppContext, ctx context.Context) {
//	c, _ := appCtx.GetPubsub().Subscribe(ctx, common.ChannelNoteCreated)
//
//	go func() {
//		for {
//			msg := <-c
//			data := msg.Data()
//
//			log.Println("RunPushNotificationAfterCreateNote:", data)
//		}
//	}()
//}

func PushNotificationAfterCreateNote(appCtx appctx.AppContext, msg *pubsub.Message) {
	data := msg.Data()
	log.Println("RunPushNotificationAfterCreateNote:", data)
}

// Who call its handler, it's an engine (Framework)

func SendNotificationAfterCreateNote(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Send notification after create note",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			log.Println("SendNotificationAfterCreateNote ", message.Data())
			return nil
		},
	}
}
