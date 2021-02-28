package notebusiness

import (
	"context"
	"fooddlv/common"
	"fooddlv/module/note/notemodel"
	"fooddlv/pubsub"
	"log"
)

type CreateNoteStore interface {
	CreateNote(context context.Context, data *notemodel.NoteCreate) error
}

type ImgStorage interface {
	ListImages(
		context context.Context,
		ids []int,
		moreKeys ...string,
	) ([]common.Image, error)
	//DeleteImages(ctx context.Context, ids []int) error
}

type createNoteBiz struct {
	store     CreateNoteStore
	imgStore  ImgStorage
	ps        pubsub.Pubsub
	requester common.Requester
}

func NewCreateNoteBiz(store CreateNoteStore, imgStore ImgStorage, ps pubsub.Pubsub, requester common.Requester) *createNoteBiz {
	return &createNoteBiz{store: store, imgStore: imgStore, ps: ps, requester: requester}
}

func (biz *createNoteBiz) CreateNewNote(ctx context.Context, data *notemodel.NoteCreate) error {
	//data.UserId = biz.requester.GetUserId()

	imgs, err := biz.imgStore.ListImages(ctx, []int{data.CoverImgId})

	if err != nil {
		return common.ErrCannotCreateEntity(notemodel.EntityName, err)
	}

	if len(imgs) == 0 {
		return common.ErrCannotCreateEntity(notemodel.EntityName, err)
	}

	data.Cover = &imgs[0]

	if err := biz.store.CreateNote(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(notemodel.EntityName, err)
	}

	go func() {
		if err := biz.ps.Publish(ctx, common.ChannelNoteCreated, pubsub.NewMessage(data)); err != nil {
			log.Println(err)
		}

		//deleteImgJob := asyncjob.NewJob(func(ctx context.Context) error {
		//	return biz.imgStore.DeleteImages(ctx, []int{data.CoverImgId})
		//})
		//
		//group := asyncjob.NewGroup(false, deleteImgJob)
		//_ = group.Run(ctx)
	}()

	return nil
}
