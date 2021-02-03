package main

import (
	"fmt"
	"fooddlv/appctx"
	"fooddlv/appctx/uploadprovider"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	dsn := os.Getenv("DB_CONNECTION_STRING")
	secret := os.Getenv("SECRET_KEY")

	s3BucketName := os.Getenv("S3_BUCKET_NAME")
	s3Region := os.Getenv("S3_REGION")
	s3APIKey := os.Getenv("S3_API_KEY")
	s3Secret := os.Getenv("S3_SECRET")
	s3Domain := fmt.Sprintf("https://%s", os.Getenv("S3_DOMAIN"))

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3Secret, s3Domain)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	appCtx := appctx.NewAppContext(db, secret, s3Provider)

	r := gin.Default()

	setupRouter(r, appCtx)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
