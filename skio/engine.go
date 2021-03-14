package skio

import (
	"context"
	"errors"
	"fmt"
	"fooddlv/appctx"
	"fooddlv/appctx/tokenprovider/jwt"
	"fooddlv/module/user/userstorage"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"log"
	"sync"
)

type RealtimeEngine interface {
	UserSockets(userId int) []AppSocket
	EmitToRoom(room string) error
	EmitToUser(userId int, key string, data interface{}) error
	Run(ctx appctx.AppContext, engine *gin.Engine) error
	//Emit(userId int) error
}

type rtEngine struct {
	server  *socketio.Server
	storage map[int][]AppSocket
	locker  *sync.RWMutex
}

func NewEngine() *rtEngine {
	return &rtEngine{
		storage: make(map[int][]AppSocket),
		locker:  new(sync.RWMutex),
	}
}

func (engine *rtEngine) saveAppSocket(userId int, appSck AppSocket) {
	engine.locker.Lock()

	if v, ok := engine.storage[userId]; ok {
		engine.storage[userId] = append(v, appSck)
	} else {
		engine.storage[userId] = []AppSocket{appSck}
	}

	engine.locker.Unlock()
}

func (engine *rtEngine) getAppSocket(userId int) []AppSocket {
	engine.locker.RLock()
	defer engine.locker.RUnlock()

	return engine.storage[userId]
}

func (engine *rtEngine) removeAppSocket(userId int, appSck AppSocket) {
	engine.locker.Lock()
	defer engine.locker.Unlock()

	if v, ok := engine.storage[userId]; ok {
		for i := range v {
			if v[i] == appSck {
				engine.storage[userId] = append(v[:i], v[i+1:]...)
				break
			}
		}
	}
}

func (engine *rtEngine) UserSockets(userId int) []AppSocket {
	var sockets []AppSocket

	if scks, ok := engine.storage[userId]; ok {
		return scks
	}

	return sockets
}

func (engine *rtEngine) EmitToRoom(room string) error {
	return nil
}

func (engine *rtEngine) EmitToUser(userId int, key string, data interface{}) error {
	sockets := engine.getAppSocket(userId)

	for _, s := range sockets {
		s.Emit(key, data)
	}

	return nil
}

func (engine *rtEngine) Run(appCtx appctx.AppContext, r *gin.Engine) error {
	server, _ := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{websocket.Default},
	})

	engine.server = server

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID(), " IP:", s.RemoteAddr(), s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	server.OnEvent("/", "authenticate", func(s socketio.Conn, token string) {
		db := appCtx.GetDBConnection()
		store := userstorage.NewSQLStore(db)

		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			s.Emit("authentication_failed", err.Error())
			s.Close()
			return
		}

		user, err := store.FindUser(context.Background(), map[string]interface{}{"id": payload.UserId})

		if err != nil {
			s.Emit("authentication_failed", err.Error())
			s.Close()
			return
		}

		if user.Status == 0 {
			s.Emit("authentication_failed", errors.New("you has been banned/deleted"))
			s.Close()
			return
		}

		appSck := NewAppSocket(s, user)
		engine.saveAppSocket(user.Id, appSck)

		s.Emit("authenticated", user)

		log.Println(s.ID(), token)
	})

	//
	//type A struct {
	//	Age int `json:"age"`
	//}
	//
	//server.OnEvent("/", "notice", func(s socketio.Conn, msg A) {
	//	fmt.Println("server receive notice:", msg.Age)
	//
	//	msg.Age = 2
	//	s.Emit("reply", msg)
	//})
	//
	//server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	//	s.SetContext(msg)
	//	return "recv " + msg
	//})
	//
	//server.OnEvent("/", "bye", func(s socketio.Conn) string {
	//	last := s.Context().(string)
	//	s.Emit("bye", last)
	//	s.Close()
	//	return last
	//})
	//
	//server.OnEvent("/", "noteSumit", func(s socketio.Conn) string {
	//	last := s.Context().(string)
	//	s.Emit("bye", last)
	//	s.Close()
	//	return last
	//})

	go server.Serve()

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	return nil
}

// Server sockets: A, B, C
// Client emit data to A, how we emit data to B, C

// A, B, C subscribe to broker pubsub
// Client emit data to A. A publish data to broker
// A, B, C receive data => Emit to data to its clients
