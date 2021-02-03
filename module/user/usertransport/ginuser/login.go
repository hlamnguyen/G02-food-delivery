package ginuser

import (
	"fooddlv/appctx"
	"fooddlv/appctx/hasher"
	"fooddlv/appctx/tokenprovider/jwt"
	"fooddlv/common"
	"fooddlv/module/user/userbusiness"
	"fooddlv/module/user/usermodel"
	"fooddlv/module/user/userstorage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {

		var loginUserData usermodel.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetDBConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()

		business := userbusiness.NewLoginBusiness(store, tokenProvider, md5, appctx.NewTokenConfig())
		account, err := business.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
