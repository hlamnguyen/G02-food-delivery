package uploadprovider

import (
	"context"
	"fooddlv/common"
	"mime/multipart"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, file multipart.File, dst string) (*common.Image, error)
}
