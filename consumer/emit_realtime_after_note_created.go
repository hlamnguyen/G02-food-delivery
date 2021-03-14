package consumer

import (
	"context"
	"fooddlv/appctx"
	"fooddlv/pubsub"
	"fooddlv/skio"
	"log"
)

func EmitRealtimeAfterCreateNote(appCtx appctx.AppContext, rtEngine skio.RealtimeEngine) consumerJob {
	return consumerJob{
		Title: "Emit realtime data after create note",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			log.Println("EmitRealtimeAfterCreateNote ", message.Data())
			_ = rtEngine.EmitToUser(1, "NoteCreated", message.Data())
			return nil
		},
	}
}
