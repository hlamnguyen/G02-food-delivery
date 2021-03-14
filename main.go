package main

import (
	"fmt"
	"fooddlv/appctx"
	"fooddlv/appctx/uploadprovider"
	"fooddlv/consumer"
	"fooddlv/pubsub/pblocal"
	"fooddlv/skio"
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

	r := gin.Default()

	realtimeEngine := skio.NewEngine()
	_ = realtimeEngine.Run(appCtx, r)

	consumerEngine := consumer.NewEngine(appCtx, realtimeEngine)
	consumerEngine.Start()

	//consumer.Setup(appCtx, context.Background())

	r.GET("/.well-known/acme-challenge/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	r.StaticFile("/demo/", "./demo.html")

	setupRouter(r, appCtx)

	je, err := jg.NewExporter(jg.Options{
		AgentEndpoint: "localhost:6831",
		Process:       jg.Process{ServiceName: "Food-Delivery"},
	})

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

	//startSocketIOServer(r, appCtx)

	http.ListenAndServe(
		":8080",
		&ochttp.Handler{
			Handler: r,
		},
	)

	//r.Run()

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Each client has a connection (web socket)
// 1. Who is this connection (need authentication)
// 2. In case server wants to emit data to specific user, we can't.
//
//func startSocketIOServer(engine *gin.Engine, appCtx appctx.AppContext) {
//	server, _ := socketio.NewServer(&engineio.Options{
//		Transports: []transport.Transport{websocket.Default},
//	})
//
//	server.OnConnect("/", func(s socketio.Conn) error {
//		s.SetContext("")
//		fmt.Println("connected:", s.ID(), " IP:", s.RemoteAddr(), s.ID())
//
//		//go func() {
//		//	i := 0
//		//	for {
//		//		i++
//		//		s.Emit("test", i)
//		//		time.Sleep(time.Second)
//		//	}
//		//}()
//		return nil
//	})
//
//	server.OnError("/", func(s socketio.Conn, e error) {
//		fmt.Println("meet error:", e)
//	})
//
//	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
//		fmt.Println("closed", reason)
//		// Remove socket from socket engine (from app context)
//	})
//
//	server.OnEvent("/", "authenticate", func(s socketio.Conn, token string) {
//		// Validate token
//		// If false: s.Close(), and return
//
//		// If true
//		// => UserId
//		// Fetch db find user by Id
//		// Here: s belongs to who? (user_id)
//		// We need a map[user_id][]socketio.Conn
//
//		db := appCtx.GetDBConnection()
//		store := userstorage.NewSQLStore(db)
//
//		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
//
//		payload, err := tokenProvider.Validate(token)
//		if err != nil {
//			s.Emit("authentication_failed", err.Error())
//			s.Close()
//			return
//		}
//
//		user, err := store.FindUser(context.Background(), map[string]interface{}{"id": payload.UserId})
//
//		if err != nil {
//			s.Emit("authentication_failed", err.Error())
//			s.Close()
//			return
//		}
//
//		if user.Status == 0 {
//			s.Emit("authentication_failed", errors.New("you has been banned/deleted"))
//			s.Close()
//			return
//		}
//
//		appSck := skio.NewAppSocket(s, user)
//
//		log.Println(s.ID(), token)
//	})
//
//	server.OnEvent("/", "test", func(s socketio.Conn, msg string) {
//		log.Println(msg)
//	})
//
//	type A struct {
//		Age int `json:"age"`
//	}
//
//	server.OnEvent("/", "notice", func(s socketio.Conn, msg A) {
//		fmt.Println("server receive notice:", msg.Age)
//
//		msg.Age = 2
//		s.Emit("reply", msg)
//	})
//
//	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
//		s.SetContext(msg)
//		return "recv " + msg
//	})
//
//	server.OnEvent("/", "bye", func(s socketio.Conn) string {
//		last := s.Context().(string)
//		s.Emit("bye", last)
//		s.Close()
//		return last
//	})
//
//	server.OnEvent("/", "noteSumit", func(s socketio.Conn) string {
//		last := s.Context().(string)
//		s.Emit("bye", last)
//		s.Close()
//		return last
//	})
//
//	go server.Serve()
//
//	engine.GET("/socket.io/*any", gin.WrapH(server))
//	engine.POST("/socket.io/*any", gin.WrapH(server))
//
//	engine.StaticFile("/demo/", "./demo.html")
//}
