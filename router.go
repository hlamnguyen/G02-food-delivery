package main

import (
	"fooddlv/appctx"
	"fooddlv/middleware"
	"fooddlv/module/note/notetransport/ginnote"
	"fooddlv/module/upload/uploadtransport/ginupload"
	"fooddlv/module/user/usertransport/ginuser"
	"github.com/gin-gonic/gin"
)

func setupRouter(r *gin.Engine, appCtx appctx.AppContext) {
	r.Use(middleware.Recover(appCtx))
	v1 := r.Group("/v1")

	v1.POST("/register", ginuser.Register(appCtx))
	v1.POST("/login", ginuser.Login(appCtx))

	v1.POST("/upload", ginupload.Upload(appCtx))
	v1.GET("/profile", middleware.RequiredAuth(appCtx), ginuser.GetProfile(appCtx))

	//users := v1.Group("users", middleware.RequiredAuth(appCtx))

	notes := v1.Group("/notes", middleware.RequiredAuth(appCtx))
	{
		notes.GET("", ginnote.ListNote(appCtx))
		notes.POST("", ginnote.CreateNote(appCtx))
		notes.GET("/:note-id", ginnote.GetNote(appCtx))
		notes.PUT("/:note-id", ginnote.UpdateNote(appCtx))
		notes.DELETE("/:note-id", ginnote.DeleteNote(appCtx))
	}

	r.Static("/upload", "./static")
}

func setupAdminRouter(r *gin.Engine, appCtx appctx.AppContext) {
	r.Use(middleware.Recover(appCtx))
	//admin := r.Group("/v1/admin")
}
