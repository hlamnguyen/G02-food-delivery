package consumer

import (
	"context"
	"fooddlv/appctx"
	"fooddlv/common"
	"fooddlv/module/upload/uploadstorage"
	"fooddlv/pubsub"
)

type HasImageIds interface {
	GetImageIds() []int
}

func RunDeleteImageRecordAfterCreateNote(appCtx appctx.AppContext, ctx context.Context) {
	c, _ := appCtx.GetPubsub().Subscribe(ctx, common.ChannelNoteCreated)

	go func() {
		for {
			msg := <-c
			if data, ok := msg.Data().(HasImageIds); ok {
				uploadstorage.NewSQLStore(appCtx.GetDBConnection()).DeleteImages(ctx, data.GetImageIds())
			}
		}
	}()
}

func DeleteImageRecordAfterCreateNote(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Delete images records after create note",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			if data, ok := message.Data().(HasImageIds); ok {
				return uploadstorage.NewSQLStore(appCtx.GetDBConnection()).DeleteImages(ctx, data.GetImageIds())
			}

			return nil
		},
	}
}
