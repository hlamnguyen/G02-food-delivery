package consumer

import (
	"context"
	"fooddlv/appctx"
	"fooddlv/common"
	"fooddlv/module/upload/uploadstorage"
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
