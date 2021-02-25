package notebusiness

import (
	"context"
	"fooddlv/common"
	"fooddlv/module/note/notemodel"
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
	DeleteImages(ctx context.Context, ids []int) error
}

type createNoteBiz struct {
	store     CreateNoteStore
	imgStore  ImgStorage
	requester common.Requester
}

func NewCreateNoteBiz(store CreateNoteStore, imgStore ImgStorage, requester common.Requester) *createNoteBiz {
	return &createNoteBiz{store: store, imgStore: imgStore, requester: requester}
}

func (biz *createNoteBiz) CreateNewNote(context context.Context, data *notemodel.NoteCreate) error {
	//data.UserId = biz.requester.GetUserId()

	imgs, err := biz.imgStore.ListImages(context, []int{data.CoverImgId})

	if err != nil {
		return common.ErrCannotCreateEntity(notemodel.EntityName, err)
	}

	if len(imgs) == 0 {
		return common.ErrCannotCreateEntity(notemodel.EntityName, err)
	}

	data.Cover = &imgs[0]

	if err := biz.store.CreateNote(context, data); err != nil {
		return common.ErrCannotCreateEntity(notemodel.EntityName, err)
	}

	go func() {
		_ = biz.imgStore.DeleteImages(context, []int{data.CoverImgId})
	}()

	return nil
}
