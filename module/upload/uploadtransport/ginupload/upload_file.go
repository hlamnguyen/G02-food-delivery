package ginupload

import (
	"fooddlv/appctx"
	"fooddlv/common"
	"fooddlv/module/upload/uploadbusiness"
	"github.com/gin-gonic/gin"
	_ "image/jpeg"
	_ "image/png"
)

func Upload(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		folder := c.DefaultPostForm("folder", "img")

		file, err := fileHeader.Open()

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		defer file.Close()

		biz := uploadbusiness.NewUploadBiz(appCtx.UploadProvider())
		img, err := biz.Upload(c.Request.Context(), file, folder, fileHeader.Filename)

		if err != nil {
			panic(err)
		}
		c.JSON(200, common.SimpleSuccessResponse(img))
	}
}
