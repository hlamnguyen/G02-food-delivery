package main

import (
	"context"
	"errors"
	"fooddlv/common/asyncjob"
	"log"
	"sync"
	"time"
)

type safeStorage struct {
	data   map[string]interface{}
	locker *sync.RWMutex
}

func (s *safeStorage) Read(key string) interface{} {
	v := s.data[key]
	return v
}

func (s *safeStorage) Write(key string, v interface{}) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.data[key] = v
}

func main() {
	myStore := &safeStorage{
		data:   make(map[string]interface{}),
		locker: new(sync.RWMutex),
	}

	j1 := asyncjob.NewJob(func(ctx context.Context) error {
		log.Println("I am job 1")
		time.Sleep(time.Second * 5)

		myStore.Write("something", "hihi")

		return nil
		//return nil
	})

	//if err := j1.Execute(context.Background()); err != nil {
	//	log.Println(err)
	//
	//	for {
	//		err = j1.Retry(context.Background())
	//		if err == nil || j1.State() == asyncjob.StateRetryFailed {
	//			break
	//		}
	//	}
	//}
	//
	//log.Println("Done job")

	j2 := asyncjob.NewJob(func(ctx context.Context) error {
		log.Println("I am job 2")
		time.Sleep(time.Second)

		//m["a"] = 1 // cause crash when running in goroutines
		//delete(m, "a") // cause crash when running in goroutines
		//for {
		//	now := time.Now().UTC()
		//	if now.Hour() == 12 && now.Minute() == 00 && now.Second() == 00 {
		//		// do something
		//	}
		//	time.Sleep(time.Second)
		//}

		//
		//ticker := time.NewTicker(time.Second)
		//for {
		//	<-ticker.C
		//}

		return errors.New("job 2 error")
	})
	//
	//j2.SetRetryDurations([]time.Duration{time.Second * 2})

	group := asyncjob.NewGroup(true, j1, j2)
	err := group.Run(context.Background())
	log.Println("Group result:", err)

	//group := asyncjob.NewGroup(true, j1, j2)
	//err := group.Run(context.Background())
	//
	//log.Println("Group result:", err)
}
