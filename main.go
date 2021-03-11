package main

import (
	"fmt"
	"fooddlv/appctx"
	"fooddlv/appctx/uploadprovider"
	"fooddlv/consumer"
	"fooddlv/pubsub/pblocal"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	jg "go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	log.Println("hello")
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

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

	ps := pblocal.NewPubSub()

	appCtx := appctx.NewAppContext(db, secret, s3Provider, ps)

	consumerEngine := consumer.NewEngine(appCtx)
	consumerEngine.Start()

	//consumer.Setup(appCtx, context.Background())

	r := gin.Default()

	r.GET("/.well-known/acme-challenge/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	setupRouter(r, appCtx)

	je, err := jg.NewExporter(jg.Options{
		AgentEndpoint: "localhost:6831",
		Process:       jg.Process{ServiceName: "Food-Delivery"},
	})

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

	http.ListenAndServe(
		":8080",
		&ochttp.Handler{
			Handler: r,
		},
	)

	//r.Run()

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
