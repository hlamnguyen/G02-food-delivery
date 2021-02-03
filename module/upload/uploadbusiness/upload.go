package uploadbusiness

import (
	"context"
	"fmt"
	"fooddlv/common"
	"fooddlv/module/upload/uploadmodel"
	"image"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, file multipart.File, dst string) (*common.Image, error)
}

type uploadBiz struct {
	provider UploadProvider
}

func NewUploadBiz(provider UploadProvider) *uploadBiz {
	return &uploadBiz{provider: provider}
}

func (biz *uploadBiz) Upload(ctx context.Context, file multipart.File, folder, fileName string) (*common.Image, error) {
	w, h, err := getImageDimension(file)

	if err != nil {
		return nil, uploadmodel.ErrFileIsNotImage(err)
	}

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	fileExt := filepath.Ext(fileName)
	fileName = fmt.Sprintf("%d%s", time.Now().Nanosecond(), fileExt)

	img, err := biz.provider.SaveFileUploaded(ctx, file, fmt.Sprintf("%s/%s", folder, fileName))

	if err != nil {
		return nil, uploadmodel.ErrCannotSaveFile(err)
	}

	img.Width = w
	img.Height = h
	img.CloudName = "s3" // should be set in provider
	img.Extension = fileExt

	return img, nil
}

func getImageDimension(reader io.Reader) (int, int, error) {
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		log.Println("err: ", err)
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}
