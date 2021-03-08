package ginupload

import (
	"fooddlv/appctx"
	"fooddlv/common"
	"fooddlv/module/upload/uploadbusiness"
	"fooddlv/module/upload/uploadstorage"
	"github.com/gin-gonic/gin"
	_ "image/jpeg"
	_ "image/png"
)

//Upload File to S3
//1. Get image/file from request header
//2. Check file is real image
//3. Save image
//1. Save to local service
//2. Save to cloud storage (S3)
//4. Improve security

func Upload(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := appCtx.GetDBConnection()

		fileHeader, err := c.FormFile("file")

		//c.SaveUploadedFile()

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		folder := c.DefaultPostForm("folder", "img")

		file, err := fileHeader.Open()

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		defer file.Close() // we can close here

		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		imgStore := uploadstorage.NewSQLStore(db)
		biz := uploadbusiness.NewUploadBiz(appCtx.UploadProvider(), imgStore)
		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)

		if err != nil {
			panic(err)
		}
		c.JSON(200, common.SimpleSuccessResponse(img.Id))
	}
}
